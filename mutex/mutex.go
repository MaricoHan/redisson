package mutex

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson/pkg/types"
	"github.com/MaricoHan/redisson/pkg/util"
)

var mutexScript = struct {
	lockScript    string
	renewalScript string
	unlockScript  string
}{}

type Mutex struct {
	root *Root
	baseMutex
}

func NewMutex(root *Root, name string, options ...Option) *Mutex {
	base := baseMutex{
		Name:   name,
		pubSub: root.Client.Subscribe(context.Background(), util.ChannelName(name)),
	}

	for i := range options {
		options[i].Apply(&base)
	}

	(&base).checkAndInit()

	return &Mutex{
		root:      root,
		baseMutex: base,
	}
}

func (m *Mutex) Lock() error {
	// 单位：ms
	expiration := int64(m.expiration / time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), m.waitTimeout)
	defer cancel()
	goID := util.GoID()
	err := m.tryLock(ctx, goID, expiration)
	if err != nil {
		return err
	}
	// 加锁成功，开个协程，定时续锁
	go func() {
		ticker := time.NewTicker(m.expiration / 3).C
		for range ticker {
			res, err := m.root.Client.Eval(context.TODO(), mutexScript.renewalScript, []string{m.Name}, expiration).Int64()
			if err != nil || res == 0 {
				return
			}
		}
	}()

	return nil
}

func (m *Mutex) tryLock(ctx context.Context, goID int64, expiration int64) error {
	// 尝试加锁
	pTTL, err := m.lockInner(goID, expiration)
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
		return m.tryLock(ctx, goID, expiration)
	case <-m.pubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		return m.tryLock(ctx, goID, expiration)
	}
}

func (m *Mutex) lockInner(goID, expiration int64) (int64, error) {
	pTTL, err := m.root.Client.Eval(context.TODO(), mutexScript.lockScript, []string{m.Name}, m.root.UUID+":"+strconv.FormatInt(goID, 10), expiration).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return pTTL.(int64), nil
}

func (m *Mutex) Unlock() error {
	goID := util.GoID()
	return m.unlockInner(goID)
}

func (m *Mutex) unlockInner(goID int64) error {
	res, err := m.root.Client.Eval(context.TODO(), mutexScript.unlockScript, []string{m.Name, util.ChannelName(m.Name)}, m.root.UUID+":"+strconv.FormatInt(goID, 10), 1).Int64()
	if err != nil {
		return err
	}
	if res == 0 {
		return types.ErrMismatch
	}

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
	return redis.call('pexpire',KEYS[1],ARGV[1])
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
