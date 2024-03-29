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

	return &Mutex{
		root:      root,
		baseMutex: base,
	}
}

func (m *Mutex) Lock(ctx context.Context) error {
	// 单位：ms
	pExpireNum := int64(m.options.expiration / time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, m.options.waitTimeout)
	defer cancel()

	var err error
	// 先订阅，再申请锁
	if m.pubSub == nil {
		m.pubSub = pubsub.Subscribe(utils.ChannelName(m.Name))
	}

	// 申请锁
	clientID := m.root.UUID + ":" + strconv.FormatInt(utils.GoID(), 10)
	if err = m.tryLock(ctx, clientID, pExpireNum); err != nil {
		return err
	}

	// 加锁成功，开个协程，定时续锁
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()

		ticker := time.NewTicker(m.options.expiration / 3)
		defer ticker.Stop()

		// 上传脚本
		if mutexScript.renewalScriptSha == "" {
			mutexScript.renewalScriptSha, err = m.root.Client.ScriptLoad(ctx, mutexScript.renewalScript).Result()
			if err != nil {
				return
			}
		}

		for {
			select {
			case <-m.release:
				return
			case <-ticker.C:
				res, err := m.root.Client.EvalSha(context.TODO(), mutexScript.renewalScriptSha, []string{m.Name}, pExpireNum, clientID).Int64()
				if err != nil || res == 0 {
					return
				}
			}
		}
	}()
	wg.Wait() // 等待协程启动成功

	return nil
}

func (m *Mutex) tryLock(ctx context.Context, clientID string, pExpireNum int64) error {
	// 尝试加锁
	pTTL, err := m.lockInner(ctx, clientID, pExpireNum)
	if err != nil {
		return err
	}
	if pTTL == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		// 申请锁的耗时如果大于等于最大等待时间，则申请锁失败.
		return types.ErrWaitTimeout
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对“redis 中存在未维护的锁”，即当锁自然过期后，并不会发布通知的锁
		return m.tryLock(ctx, clientID, pExpireNum)
	case <-m.pubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		return m.tryLock(ctx, clientID, pExpireNum)
	}
}

func (m *Mutex) lockInner(ctx context.Context, clientID string, pExpireNum int64) (int64, error) {
	// 上传脚本
	if mutexScript.lockScriptSha == "" {
		var err error
		mutexScript.lockScriptSha, err = m.root.Client.ScriptLoad(ctx, mutexScript.lockScript).Result()
		if err != nil {
			return 0, fmt.Errorf("load lock script err: %w", err)
		}
	}

	pTTL, err := m.root.Client.EvalSha(ctx, mutexScript.lockScriptSha, []string{m.Name}, clientID, pExpireNum).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return pTTL.(int64), nil
}

func (m *Mutex) Unlock(ctx context.Context) error {
	goID := utils.GoID()

	if err := m.unlockInner(ctx, goID); err != nil {
		return fmt.Errorf("unlock err: %w", err)
	}

	return nil
}

func (m *Mutex) unlockInner(ctx context.Context, goID int64) error {
	// 上传脚本
	if mutexScript.unlockScriptSha == "" {
		var err error
		mutexScript.unlockScriptSha, err = m.root.Client.ScriptLoad(ctx, mutexScript.unlockScript).Result()
		if err != nil {
			return fmt.Errorf("load unlock script err: %w", err)
		}
	}

	res, err := m.root.Client.EvalSha(
		ctx,
		mutexScript.unlockScriptSha,
		[]string{m.Name, m.root.RedisChannelName},
		m.root.UUID+":"+strconv.FormatInt(goID, 10),
		m.Name+":unlock",
	).Int64()
	if err != nil {
		return err
	}
	if res == 0 {
		return types.ErrMismatch
	}

	// 释放资源
	m.pubSub.Close()
	close(m.release) // 通知续锁协程退出

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
