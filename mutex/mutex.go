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

var mutexScript = struct {
	lockScript    string
	lockScriptSha string

	renewalScript    string
	renewalScriptSha string

	unlockScript    string
	unlockScriptSha string
}{}

type Mutex struct {
	root *Root
	*baseMutex
}

func NewMutex(root *Root, name string, opts ...Option) *Mutex {
	base := &baseMutex{
		Name:    name,
		release: make(chan struct{}),
		options: &options{},
	}
	for i := range opts {
		opts[i](base.options)
	}

	base.options.checkAndInit()

	root.Logger.Debugf("创建互斥锁实例: %s, 过期时间: %v, 等待超时: %v",
		name, base.options.expiration, base.options.waitTimeout)

	return &Mutex{
		root:      root,
		baseMutex: base,
	}
}

func (m *Mutex) Lock(ctx context.Context) error {
	// 单位：ms
	pExpireNum := int64(m.options.expiration / time.Millisecond)

	m.root.Logger.Debugf("尝试获取互斥锁: %s, 过期时间: %dms", m.Name, pExpireNum)

	ctx, cancel := context.WithTimeout(ctx, m.options.waitTimeout)
	defer cancel()

	var err error
	// 先订阅，再申请锁
	if m.pubSub == nil {
		m.pubSub = pubsub.Subscribe(utils.ChannelName(m.Name))
		m.root.Logger.Debugf("订阅锁通道: %s", utils.ChannelName(m.Name))
	}

	// 申请锁
	clientID := m.root.UUID + ":" + strconv.FormatInt(utils.GoID(), 10)
	if err = m.tryLock(ctx, clientID, pExpireNum); err != nil {
		m.root.Logger.Errorf("获取互斥锁失败: %s, 客户端ID: %s, 错误: %v", m.Name, clientID, err)
		return err
	}

	m.root.Logger.Infof("成功获取互斥锁: %s, 客户端ID: %s", m.Name, clientID)

	// 加锁成功，开个协程，定时续锁
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()

		ticker := time.NewTicker(m.options.expiration / 3)
		defer ticker.Stop()

		m.root.Logger.Debugf("启动互斥锁续期协程: %s, 续期间隔: %v", m.Name, m.options.expiration/3)

		// 上传脚本
		if mutexScript.renewalScriptSha == "" {
			mutexScript.renewalScriptSha, err = m.root.Client.ScriptLoad(ctx, mutexScript.renewalScript).Result()
			if err != nil {
				m.root.Logger.Errorf("加载互斥锁续期脚本失败: %v", err)
				return
			}
			m.root.Logger.Debugf("加载互斥锁续期脚本成功: %s", mutexScript.renewalScriptSha)
		}

		for {
			select {
			case <-m.release:
				m.root.Logger.Debugf("互斥锁续期协程收到退出信号: %s", m.Name)
				return
			case <-ticker.C:
				res, err := m.root.Client.EvalSha(context.TODO(), mutexScript.renewalScriptSha, []string{m.Name}, pExpireNum, clientID).Int64()
				if err != nil {
					m.root.Logger.Errorf("互斥锁续期失败: %s, 错误: %v", m.Name, err)
					return
				}
				if res == 0 {
					m.root.Logger.Warnf("互斥锁续期失败，锁已不存在或已被其他客户端获取: %s", m.Name)
					return
				}
				m.root.Logger.Debugf("互斥锁续期成功: %s", m.Name)
			}
		}
	}()
	wg.Wait() // 等待协程启动成功

	return nil
}

func (m *Mutex) tryLock(ctx context.Context, clientID string, pExpireNum int64) error {
	// 尝试加锁
	m.root.Logger.Debugf("尝试获取互斥锁: %s, 客户端ID: %s", m.Name, clientID)
	pTTL, err := m.lockInner(ctx, clientID, pExpireNum)
	if err != nil {
		m.root.Logger.Errorf("获取互斥锁内部操作失败: %s, 错误: %v", m.Name, err)
		return err
	}
	if pTTL == 0 {
		m.root.Logger.Debugf("成功获取互斥锁: %s", m.Name)
		return nil
	}

	m.root.Logger.Debugf("互斥锁已被占用: %s, TTL: %dms, 等待解锁或过期", m.Name, pTTL)

	select {
	case <-ctx.Done():
		// 申请锁的耗时如果大于等于最大等待时间，则申请锁失败.
		m.root.Logger.Warnf("获取互斥锁等待超时: %s", m.Name)
		return types.ErrWaitTimeout
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对"redis 中存在未维护的锁"，即当锁自然过期后，并不会发布通知的锁
		m.root.Logger.Debugf("互斥锁等待过期后重试: %s", m.Name)
		return m.tryLock(ctx, clientID, pExpireNum)
	case <-m.pubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		m.root.Logger.Debugf("收到互斥锁解锁通知，尝试获取: %s", m.Name)
		return m.tryLock(ctx, clientID, pExpireNum)
	}
}

func (m *Mutex) lockInner(ctx context.Context, clientID string, pExpireNum int64) (int64, error) {
	// 上传脚本
	if mutexScript.lockScriptSha == "" {
		var err error
		m.root.Logger.Debugf("加载互斥锁获取脚本")
		mutexScript.lockScriptSha, err = m.root.Client.ScriptLoad(ctx, mutexScript.lockScript).Result()
		if err != nil {
			m.root.Logger.Errorf("加载互斥锁获取脚本失败: %v", err)
			return 0, fmt.Errorf("load lock script err: %w", err)
		}
		m.root.Logger.Debugf("加载互斥锁获取脚本成功: %s", mutexScript.lockScriptSha)
	}

	pTTL, err := m.root.Client.EvalSha(ctx, mutexScript.lockScriptSha, []string{m.Name}, clientID, pExpireNum).Result()
	if err == redis.Nil {
		m.root.Logger.Debugf("互斥锁获取成功: %s", m.Name)
		return 0, nil
	}

	if err != nil {
		m.root.Logger.Errorf("执行互斥锁获取脚本失败: %v", err)
		return 0, err
	}

	return pTTL.(int64), nil
}

func (m *Mutex) Unlock(ctx context.Context) error {
	goID := utils.GoID()
	clientID := m.root.UUID + ":" + strconv.FormatInt(goID, 10)

	m.root.Logger.Debugf("尝试释放互斥锁: %s, 客户端ID: %s", m.Name, clientID)

	if err := m.unlockInner(ctx, goID); err != nil {
		m.root.Logger.Errorf("释放互斥锁失败: %s, 错误: %v", m.Name, err)
		return fmt.Errorf("unlock err: %w", err)
	}

	m.root.Logger.Infof("成功释放互斥锁: %s", m.Name)
	return nil
}

func (m *Mutex) unlockInner(ctx context.Context, goID int64) error {
	clientID := m.root.UUID + ":" + strconv.FormatInt(goID, 10)

	// 上传脚本
	if mutexScript.unlockScriptSha == "" {
		var err error
		m.root.Logger.Debugf("加载互斥锁释放脚本")
		mutexScript.unlockScriptSha, err = m.root.Client.ScriptLoad(ctx, mutexScript.unlockScript).Result()
		if err != nil {
			m.root.Logger.Errorf("加载互斥锁释放脚本失败: %v", err)
			return fmt.Errorf("load unlock script err: %w", err)
		}
		m.root.Logger.Debugf("加载互斥锁释放脚本成功: %s", mutexScript.unlockScriptSha)
	}

	res, err := m.root.Client.EvalSha(
		ctx,
		mutexScript.unlockScriptSha,
		[]string{m.Name, m.root.RedisChannelName},
		clientID,
		m.Name+":unlock",
	).Int64()
	if err != nil {
		m.root.Logger.Errorf("执行互斥锁释放脚本失败: %v", err)
		return err
	}
	if res == 0 {
		m.root.Logger.Warnf("互斥锁释放失败，锁不存在或不匹配: %s, 客户端ID: %s", m.Name, clientID)
		return types.ErrMismatch
	}

	// 释放资源
	m.pubSub.Close()
	close(m.release) // 通知续锁协程退出
	m.root.Logger.Debugf("关闭互斥锁相关资源: %s", m.Name)

	return nil
}

func init() {
	mutexScript.lockScript = `
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

	mutexScript.renewalScript = `
	-- KEYS[1] 锁名
	-- ARGV[1] 过期时间
	-- ARGV[2] 客户端协程唯一标识
	if redis.call('get',KEYS[1])==ARGV[2] then
		return redis.call('pexpire',KEYS[1],ARGV[1])
	end
	return 0
`

	mutexScript.unlockScript = `
	-- KEYS[1] 锁名
	-- KEYS[2] 发布订阅的channel
	-- ARGV[1] 协程唯一标识：客户端标识+协程ID
	-- ARGV[2] 解锁时发布的消息
	if redis.call('exists',KEYS[1]) == 1 then
		if (redis.call('get',KEYS[1]) == ARGV[1]) then
			redis.call('del',KEYS[1])
		else
			return 0
		end
	end
	redis.call('publish',KEYS[2],ARGV[2])
	return 1
`
}
