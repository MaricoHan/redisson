package mutex

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson/pkg/types"
	"github.com/MaricoHan/redisson/pkg/utils"
	"github.com/MaricoHan/redisson/pkg/utils/pubsub"
)

var (
	rwMutexScript = struct {
		lockScript    string
		lockScriptSha string

		rLockScript    string
		rLockScriptSha string

		renewalScript    string
		renewalScriptSha string

		unlockScript    string
		unlockScriptSha string
	}{}
)

type RWMutex struct {
	root *Root
	*baseMutex
}

func NewRWMutex(r *Root, name string, opts ...Option) *RWMutex {
	base := &baseMutex{
		Name:    name,
		release: make(chan struct{}),
		options: &options{},
	}
	for i := range opts {
		opts[i](base.options)
	}

	base.options.checkAndInit()

	r.Logger.Debugf("创建读写锁实例: %s, 过期时间: %v, 等待超时: %v",
		name, base.options.expiration, base.options.waitTimeout)

	return &RWMutex{
		root:      r,
		baseMutex: base,
	}
}

func (r *RWMutex) Lock(ctx context.Context) error {
	// 单位：ms
	expiration := int64(r.options.expiration / time.Millisecond)

	r.root.Logger.Debugf("尝试获取写锁: %s, 过期时间: %dms", r.Name, expiration)

	ctx, cancel := context.WithTimeout(ctx, r.options.waitTimeout)
	defer cancel()

	var err error
	// 先订阅，再申请锁
	if r.pubSub == nil {
		r.pubSub = pubsub.Subscribe(utils.ChannelName(r.Name))
		r.root.Logger.Debugf("订阅锁通道: %s", utils.ChannelName(r.Name))
	}

	clientID := r.root.UUID + ":" + strconv.FormatInt(utils.GoID(), 10)
	if err = r.tryLock(ctx, clientID, expiration); err != nil {
		r.root.Logger.Errorf("获取写锁失败: %s, 客户端ID: %s, 错误: %v", r.Name, clientID, err)
		return err
	}

	r.root.Logger.Infof("成功获取写锁: %s, 客户端ID: %s", r.Name, clientID)

	// 加锁成功，开个协程，定时续锁
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()

		ticker := time.NewTicker(r.options.expiration / 3)
		defer ticker.Stop()

		r.root.Logger.Debugf("启动写锁续期协程: %s, 续期间隔: %v", r.Name, r.options.expiration/3)

		// 上传脚本
		if rwMutexScript.renewalScriptSha == "" {
			rwMutexScript.renewalScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.renewalScript).Result()
			if err != nil {
				r.root.Logger.Errorf("加载写锁续期脚本失败: %v", err)
				return
			}
			r.root.Logger.Debugf("加载写锁续期脚本成功: %s", rwMutexScript.renewalScriptSha)
		}

		for {
			select {
			case <-r.release:
				r.root.Logger.Debugf("写锁续期协程收到退出信号: %s", r.Name)
				return
			case <-ticker.C:
				res, err := r.root.Client.EvalSha(context.TODO(), rwMutexScript.renewalScriptSha, []string{r.Name}, expiration, clientID).Int64()
				if err != nil {
					r.root.Logger.Errorf("写锁续期失败: %s, 错误: %v", r.Name, err)
					return
				}
				if res == 0 {
					r.root.Logger.Warnf("写锁续期失败，锁已不存在或已被其他客户端获取: %s", r.Name)
					return
				}
				r.root.Logger.Debugf("写锁续期成功: %s", r.Name)
			}
		}
	}()
	wg.Wait() // 等待协程启动成功

	return nil
}

func (r *RWMutex) tryLock(ctx context.Context, clientID string, expiration int64) error {
	// 尝试加锁
	r.root.Logger.Debugf("尝试获取写锁: %s, 客户端ID: %s", r.Name, clientID)
	pTTL, err := r.lockInner(ctx, clientID, expiration)
	if err != nil {
		r.root.Logger.Errorf("获取写锁内部操作失败: %s, 错误: %v", r.Name, err)
		return err
	}
	if pTTL == 0 {
		r.root.Logger.Debugf("成功获取写锁: %s", r.Name)
		return nil
	}

	r.root.Logger.Debugf("写锁已被占用: %s, TTL: %dms, 等待解锁或过期", r.Name, pTTL)

	select {
	case <-ctx.Done():
		// 申请锁的耗时如果大于等于最大等待时间，则申请锁失败.
		r.root.Logger.Warnf("获取写锁等待超时: %s", r.Name)
		return types.ErrWaitTimeout
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对"redis 中存在未维护的锁"，即当锁自然过期后，并不会发布通知的锁
		r.root.Logger.Debugf("写锁等待过期后重试: %s", r.Name)
		return r.tryLock(ctx, clientID, expiration)
	case <-r.pubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		r.root.Logger.Debugf("收到写锁解锁通知，尝试获取: %s", r.Name)
		return r.tryLock(ctx, clientID, expiration)
	}
}

func (r *RWMutex) lockInner(ctx context.Context, clientID string, expiration int64) (int64, error) {
	// 上传脚本
	if rwMutexScript.lockScriptSha == "" {
		var err error
		r.root.Logger.Debugf("加载写锁获取脚本")
		rwMutexScript.lockScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.lockScript).Result()
		if err != nil {
			r.root.Logger.Errorf("加载写锁获取脚本失败: %v", err)
			return 0, fmt.Errorf("load lock script err: %w", err)
		}
		r.root.Logger.Debugf("加载写锁获取脚本成功: %s", rwMutexScript.lockScriptSha)
	}

	pTTL, err := r.root.Client.EvalSha(ctx, rwMutexScript.lockScriptSha, []string{r.Name}, clientID, expiration).Result()
	if err == redis.Nil {
		r.root.Logger.Debugf("写锁获取成功: %s", r.Name)
		return 0, nil
	}

	if err != nil {
		r.root.Logger.Errorf("执行写锁获取脚本失败: %v", err)
		return 0, err
	}

	return pTTL.(int64), nil
}

func (r *RWMutex) RLock(ctx context.Context) error {
	// 单位：ms
	pExpireNum := int64(r.options.expiration / time.Millisecond)

	r.root.Logger.Debugf("尝试获取读锁: %s, 过期时间: %dms", r.Name, pExpireNum)

	ctx, cancel := context.WithTimeout(ctx, r.options.waitTimeout)
	defer cancel()

	var err error
	// 先订阅，再申请锁
	if r.pubSub == nil {
		r.pubSub = pubsub.Subscribe(utils.ChannelName(r.Name))
		r.root.Logger.Debugf("订阅锁通道: %s", utils.ChannelName(r.Name))
	}

	clientID := r.root.UUID + ":" + strconv.FormatInt(utils.GoID(), 10)
	if err = r.tryRLock(ctx, clientID, pExpireNum); err != nil {
		r.root.Logger.Errorf("获取读锁失败: %s, 客户端ID: %s, 错误: %v", r.Name, clientID, err)
		return err
	}

	r.root.Logger.Infof("成功获取读锁: %s, 客户端ID: %s", r.Name, clientID)

	// 加锁成功，开个协程，定时续锁
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()

		ticker := time.NewTicker(r.options.expiration / 3)
		defer ticker.Stop()

		r.root.Logger.Debugf("启动读锁续期协程: %s, 续期间隔: %v", r.Name, r.options.expiration/3)

		// 上传脚本
		if rwMutexScript.renewalScriptSha == "" {
			rwMutexScript.renewalScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.renewalScript).Result()
			if err != nil {
				r.root.Logger.Errorf("加载读锁续期脚本失败: %v", err)
				return
			}
			r.root.Logger.Debugf("加载读锁续期脚本成功: %s", rwMutexScript.renewalScriptSha)
		}

		for {
			select {
			case <-r.release:
				r.root.Logger.Debugf("读锁续期协程收到退出信号: %s", r.Name)
				return
			case <-ticker.C:
				res, err := r.root.Client.EvalSha(context.TODO(), rwMutexScript.renewalScriptSha, []string{r.Name}, pExpireNum, clientID).Int64()
				if err != nil {
					r.root.Logger.Errorf("读锁续期失败: %s, 错误: %v", r.Name, err)
					return
				}
				if res == 0 {
					r.root.Logger.Warnf("读锁续期失败，锁已不存在或已被其他客户端获取: %s", r.Name)
					return
				}
				r.root.Logger.Debugf("读锁续期成功: %s", r.Name)
			}
		}
	}()
	wg.Wait() // 等待协程启动成功

	return nil
}

func (r *RWMutex) tryRLock(ctx context.Context, clientID string, pExpireNum int64) error {
	// 尝试加锁
	r.root.Logger.Debugf("尝试获取读锁: %s, 客户端ID: %s", r.Name, clientID)
	pTTL, err := r.rLockInner(ctx, clientID, pExpireNum)
	if err != nil {
		r.root.Logger.Errorf("获取读锁内部操作失败: %s, 错误: %v", r.Name, err)
		return err
	}
	if pTTL == 0 {
		r.root.Logger.Debugf("成功获取读锁: %s", r.Name)
		return nil
	}

	r.root.Logger.Debugf("读锁已被占用: %s, TTL: %dms, 等待解锁或过期", r.Name, pTTL)

	select {
	case <-ctx.Done():
		// 申请锁的耗时如果大于等于最大等待时间，则申请锁失败.
		r.root.Logger.Warnf("获取读锁等待超时: %s", r.Name)
		return types.ErrWaitTimeout
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对"redis 中存在未维护的锁"，即当锁自然过期后，并不会发布通知的锁
		r.root.Logger.Debugf("读锁等待过期后重试: %s", r.Name)
		return r.tryRLock(ctx, clientID, pExpireNum)
	case <-r.pubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		r.root.Logger.Debugf("收到读锁解锁通知，尝试获取: %s", r.Name)
		return r.tryRLock(ctx, clientID, pExpireNum)
	}
}

func (r *RWMutex) rLockInner(ctx context.Context, clientID string, pExpireNum int64) (int64, error) {
	// 上传脚本
	if rwMutexScript.rLockScriptSha == "" {
		var err error
		r.root.Logger.Debugf("加载读锁获取脚本")
		rwMutexScript.rLockScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.rLockScript).Result()
		if err != nil {
			r.root.Logger.Errorf("加载读锁获取脚本失败: %v", err)
			return 0, fmt.Errorf("load rlock script err: %w", err)
		}
		r.root.Logger.Debugf("加载读锁获取脚本成功: %s", rwMutexScript.rLockScriptSha)
	}

	pTTL, err := r.root.Client.EvalSha(ctx, rwMutexScript.rLockScriptSha, []string{r.Name}, clientID, pExpireNum).Result()
	if err == redis.Nil {
		r.root.Logger.Debugf("读锁获取成功: %s", r.Name)
		return 0, nil
	}

	if err != nil {
		r.root.Logger.Errorf("执行读锁获取脚本失败: %v", err)
		return 0, err
	}

	return pTTL.(int64), nil
}

func (r *RWMutex) Unlock(ctx context.Context) error {
	goID := utils.GoID()
	clientID := r.root.UUID + ":" + strconv.FormatInt(goID, 10)

	r.root.Logger.Debugf("尝试释放锁: %s, 客户端ID: %s", r.Name, clientID)

	if err := r.unlockInner(ctx, goID); err != nil {
		r.root.Logger.Errorf("释放锁失败: %s, 错误: %v", r.Name, err)
		return fmt.Errorf("unlock err: %w", err)
	}

	r.root.Logger.Infof("成功释放锁: %s", r.Name)
	return nil
}

func (r *RWMutex) unlockInner(ctx context.Context, goID int64) error {
	clientID := r.root.UUID + ":" + strconv.FormatInt(goID, 10)

	// 上传脚本
	if rwMutexScript.unlockScriptSha == "" {
		var err error
		r.root.Logger.Debugf("加载锁释放脚本")
		rwMutexScript.unlockScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.unlockScript).Result()
		if err != nil {
			r.root.Logger.Errorf("加载锁释放脚本失败: %v", err)
			return fmt.Errorf("load unlock script err: %w", err)
		}
		r.root.Logger.Debugf("加载锁释放脚本成功: %s", rwMutexScript.unlockScriptSha)
	}

	res, err := r.root.Client.EvalSha(
		ctx,
		rwMutexScript.unlockScriptSha,
		[]string{r.Name, r.root.RedisChannelName},
		clientID,
		r.Name+":unlock",
	).Int64()
	if err != nil {
		r.root.Logger.Errorf("执行锁释放脚本失败: %v", err)
		return err
	}
	if res == 0 {
		r.root.Logger.Warnf("锁释放失败，锁不存在或不匹配: %s, 客户端ID: %s", r.Name, clientID)
		return types.ErrMismatch
	}

	// 释放资源
	r.pubSub.Close()
	close(r.release) // 通知续锁协程退出
	r.root.Logger.Debugf("关闭锁相关资源: %s", r.Name)

	return nil
}

func init() {
	rwMutexScript.lockScript = `
	-- KEYS[1] 锁名
	-- ARGV[1] 协程唯一标识：客户端标识+协程ID
	-- ARGV[2] 过期时间
	if redis.call('exists',KEYS[1]) == 0 then
		redis.call('set',KEYS[1],ARGV[1])
		redis.call('pexpire',KEYS[1],ARGV[2])
		return nil
	end
	return redis.call('pttl',KEYS[1])
`

	rwMutexScript.rLockScript = `
	-- KEYS[1] 锁名
	-- ARGV[1] 协程唯一标识：客户端标识+协程ID
	-- ARGV[2] 过期时间
	local t = redis.call('type',KEYS[1])["ok"]
	if t == "string" then
		return redis.call('pttl',KEYS[1])
	else
		redis.call('hincrby',KEYS[1],ARGV[1],1)
		redis.call('pexpire',KEYS[1],ARGV[2])
		return nil
	end
`
	rwMutexScript.renewalScript = `
	-- KEYS[1] 锁名
	-- ARGV[1] 过期时间
	-- ARGV[2] 客户端协程唯一标识
	local t = redis.call('type',KEYS[1])["ok"]
	if t =="string" then
		if redis.call('get',KEYS[1])==ARGV[2] then
			return redis.call('pexpire',KEYS[1],ARGV[1])
		end
		return 0
	elseif t == "hash" then
		if redis.call('hexists',KEYS[1],ARGV[2])==0 then
			return 0
		end
		return redis.call('pexpire',KEYS[1],ARGV[1])
	else
		return 0
	end
`

	rwMutexScript.unlockScript = `
	-- KEYS[1] 锁名
	-- KEYS[2] 发布订阅的channel
	-- ARGV[1] 协程唯一标识：客户端标识+协程ID
	-- ARGV[2] 解锁时发布的消息
	-- 返回值：0-未解锁 1-解锁且整个rw锁已被删除 2-解锁且还有其他r锁存在
	local t = redis.call('type',KEYS[1])["ok"]
	if  t == "hash" then
		if redis.call('hexists',KEYS[1],ARGV[1]) == 0 then
			return 0
		end
		if redis.call('hincrby',KEYS[1],ARGV[1],-1) <= 0 then
			redis.call('hdel',KEYS[1],ARGV[1])
			if (redis.call('hlen',KEYS[1]) > 0 )then
				return 2
			end
			redis.call('del',KEYS[1])
			redis.call('publish',KEYS[2],ARGV[2])
			return 1
		else
			return 2
		end
	elseif t == "none" then
			redis.call('publish',KEYS[2],ARGV[2])
			return 1
	elseif redis.call('get',KEYS[1]) == ARGV[1] then
			redis.call('del',KEYS[1])
			redis.call('publish',KEYS[2],ARGV[2])
			return 1
	else
		return 0
	end
`
}
