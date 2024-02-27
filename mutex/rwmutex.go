package mutex

import (
	"context"
	"fmt"
	"strconv"
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
		pubSub:  pubsub.Subscribe(utils.ChannelName(name)),
		release: make(chan struct{}),
		options: &options{},
	}
	for i := range opts {
		opts[i](base.options)
	}

	base.options.checkAndInit()

	return &RWMutex{
		root:      r,
		baseMutex: base,
	}
}

func (r RWMutex) Lock(ctx context.Context) error {
	// 单位：ms
	expiration := int64(r.options.expiration / time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, r.options.waitTimeout)
	defer cancel()

	clientID := r.root.UUID + ":" + strconv.FormatInt(utils.GoID(), 10)
	err := r.tryLock(ctx, clientID, expiration)
	if err != nil {
		return err
	}
	// 加锁成功，开个协程，定时续锁
	go func() {
		ticker := time.NewTicker(r.options.expiration / 3)
		defer ticker.Stop()

		// 上传脚本
		if rwMutexScript.renewalScriptSha == "" {
			rwMutexScript.renewalScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.renewalScript).Result()
			if err != nil {
				return
			}
		}

		for {
			select {
			case <-r.release:
				return
			case <-ticker.C:
				res, err := r.root.Client.EvalSha(context.TODO(), rwMutexScript.renewalScriptSha, []string{r.Name}, expiration, clientID).Int64()
				if err != nil || res == 0 {
					return
				}
			}
		}
	}()

	return nil

}

func (r RWMutex) tryLock(ctx context.Context, clientID string, expiration int64) error {
	pTTL, err := r.lockInner(ctx, clientID, expiration)
	if err != nil {
		return err
	}

	if pTTL == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return types.ErrWaitTimeout
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对“redis 中存在未维护的锁”，即当锁自然过期后，并不会发布通知的锁
		return r.tryLock(ctx, clientID, expiration)
	case <-r.pubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		return r.tryLock(ctx, clientID, expiration)
	}

}

func (r RWMutex) lockInner(ctx context.Context, clientID string, expiration int64) (int64, error) {
	// 上传脚本
	if rwMutexScript.lockScriptSha == "" {
		var err error
		rwMutexScript.lockScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.lockScript).Result()
		if err != nil {
			return 0, fmt.Errorf("script load err: %w", err)
		}
	}

	pTTL, err := r.root.Client.EvalSha(ctx, rwMutexScript.lockScriptSha, []string{r.Name}, clientID, expiration).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return pTTL.(int64), nil
}

func (r RWMutex) RLock(ctx context.Context) error {
	// 单位：ms
	pExpireNum := int64(r.options.expiration / time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, r.options.waitTimeout)
	defer cancel()

	clientID := r.root.UUID + ":" + strconv.FormatInt(utils.GoID(), 10)
	err := r.tryRLock(ctx, clientID, pExpireNum)
	if err != nil {
		return err
	}
	// 加锁成功，开个协程，定时续锁
	go func() {
		ticker := time.NewTicker(r.options.expiration / 3)
		defer ticker.Stop()

		// 上传脚本
		if rwMutexScript.renewalScriptSha == "" {
			rwMutexScript.renewalScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.renewalScript).Result()
			if err != nil {
				return
			}
		}

		for {
			select {
			case <-r.release:
				return
			case <-ticker.C:
				res, err := r.root.Client.EvalSha(context.TODO(), rwMutexScript.renewalScriptSha, []string{r.Name}, pExpireNum, clientID).Int64()
				if err != nil || res == 0 {
					return
				}
			}
		}
	}()

	return nil

}

func (r RWMutex) tryRLock(ctx context.Context, clientID string, pExpireNum int64) error {
	pTTL, err := r.rLockInner(ctx, clientID, pExpireNum)
	if err != nil {
		return err
	}

	if pTTL == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return types.ErrWaitTimeout
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对“redis 中存在未维护的锁”，即当锁自然过期后，并不会发布通知的锁
		return r.tryRLock(ctx, clientID, pExpireNum)
	case <-r.pubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		return r.tryRLock(ctx, clientID, pExpireNum)
	}

}

func (r RWMutex) rLockInner(ctx context.Context, clientID string, pExpireNum int64) (int64, error) {
	// 上传脚本
	if rwMutexScript.rLockScriptSha == "" {
		var err error
		rwMutexScript.rLockScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.rLockScript).Result()
		if err != nil {
			return 0, fmt.Errorf("script load err: %w", err)
		}
	}

	pTTL, err := r.root.Client.EvalSha(ctx, rwMutexScript.rLockScriptSha, []string{r.Name}, clientID, pExpireNum).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return pTTL.(int64), nil
}
func (r RWMutex) Unlock(ctx context.Context) error {
	goID := utils.GoID()
	if err := r.unlockInner(ctx, goID); err != nil {
		return fmt.Errorf("unlock err: %w", err)
	}

	return nil
}
func (r RWMutex) unlockInner(ctx context.Context, goID int64) error {
	// 上传脚本
	if rwMutexScript.unlockScriptSha == "" {
		var err error
		rwMutexScript.unlockScriptSha, err = r.root.Client.ScriptLoad(ctx, rwMutexScript.unlockScript).Result()
		if err != nil {
			return fmt.Errorf("script load err: %w", err)
		}
	}

	res, err := r.root.Client.EvalSha(
		ctx,
		rwMutexScript.unlockScriptSha,
		[]string{r.Name, r.root.RedisChannelName},
		r.root.UUID+":"+strconv.FormatInt(goID, 10),
		r.Name+":unlock",
	).Int64()
	if err != nil {
		return err
	}
	if res == 0 {
		return types.ErrMismatch
	}

	// redis 中锁已经被删除，释放资源
	if res == 1 {
		r.pubSub.Close()
		close(r.release)
	}

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
