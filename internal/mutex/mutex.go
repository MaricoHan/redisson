package mutex

import (
	"context"
	"github.com/MaricoHan/redisson/pkg/base"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson/internal/root"
	"github.com/MaricoHan/redisson/pkg/types"
)

var (
	lockScript    string
	unlockScript  string
	renewalScript string
)

type Mutex struct {
	*root.Root
	root.BaseMutex
}

func NewMutex(r *root.Root, name string, options ...Option) *Mutex {
	baseMutex := root.BaseMutex{
		Name:   name,
		PubSub: r.Client.Subscribe(context.Background(), base.ChannelName(name)),
	}

	m := &Mutex{
		Root:      r,
		BaseMutex: baseMutex,
	}

	for i := range options {
		options[i].Apply(m)
	}

	return m.CheckAndInit()
}

func (m *Mutex) CheckAndInit() *Mutex {
	m.BaseMutex.CheckAndInit()

	return m
}

func (m *Mutex) Lock() error {
	// 单位：ms
	expiration := int64(m.Expiration / time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), m.WaitTimeout)
	defer cancel()
	goID := base.GoID()
	err := m.tryLock(ctx, goID, expiration)
	if err != nil {
		return err
	}
	// 加锁成功，开个协程，定时续锁
	go func() {
		ticker := time.NewTicker(m.Expiration / 3).C
		for range ticker {
			res, err := m.Client.Eval(context.TODO(), renewalScript, []string{m.Name}, expiration).Int64()
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
	case <-m.PubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		return m.tryLock(ctx, goID, expiration)
	}
}

func (m *Mutex) lockInner(goID, expiration int64) (int64, error) {
	pTTL, err := m.Client.Eval(context.TODO(), lockScript, []string{m.Name}, m.UUID+":"+strconv.FormatInt(goID, 10), expiration).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return pTTL.(int64), nil
}

func (m *Mutex) Unlock() error {
	goID := base.GoID()
	return m.unlockInner(goID)
}

func (m *Mutex) unlockInner(goID int64) error {
	res, err := m.Client.Eval(context.TODO(), unlockScript, []string{m.Name, base.ChannelName(m.Name)}, m.UUID+":"+strconv.FormatInt(goID, 10), 1).Int64()
	if err != nil {
		return err
	}
	if res == 0 {
		return types.ErrMismatch
	}

	return nil
}

func init() {
	path, _ := filepath.Abs(os.Args[1])
	index := strings.LastIndex(path, "/redisson")
	path = path[:index+9] + "/internal/mutex/lua"

	file, err := ioutil.ReadFile(path + "/lock.lua")
	if err != nil {
		panic(err)
	}
	lockScript = string(file)

	file, err = ioutil.ReadFile(path + "/unlock.lua")
	if err != nil {
		panic(err)
	}
	unlockScript = string(file)

	file, err = ioutil.ReadFile(path + "/renewal.lua")
	if err != nil {
		panic(err)
	}
	renewalScript = string(file)
}
