package mutex

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson/base"
	"github.com/MaricoHan/redisson/internal/root"
)

var (
	lockScript    string
	unlockScript  string
	renewalScript string
)

type Mutex struct {
	*root.Root
	Name       string
	Expiration time.Duration
	*redis.PubSub
}

func (m *Mutex) Init() *Mutex {
	if m.Expiration <= 0 {
		m.Expiration = 10 * time.Second
	}
	return m
}

func (m *Mutex) Lock() error {
	// 单位：ms
	expiration := int64(m.Expiration / time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), m.LockTimeout)
	defer cancel()
	goID := base.GoID()
	_, err := m.tryLock(ctx, goID, expiration)
	if err != nil {
		return err
	}
	// 加锁成功，开个协程，定时续锁
	go func() {
		ticker := time.Tick(m.Expiration / 3)
		for range ticker {
			res, err := m.Client.Eval(context.TODO(), renewalScript, []string{m.Name}, expiration).Int64()
			if err != nil || res == 0 {
				return
			}
		}

	}()

	return nil
}

func (m *Mutex) tryLock(ctx context.Context, goID int64, expiration int64) (bool, error) {
	// 尝试加锁
	pTTL, err := m.lockInner(goID, expiration)
	if err != nil {
		return false, err
	}
	if pTTL == 0 {
		return true, nil
	}

	select {
	case <-ctx.Done():
		// 申请锁的耗时如果大于等于最大等待时间，则申请锁失败.
		return false, nil
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对“redis 中存在未维护的锁”，即当锁自然过期后，并不会发布通知的锁
		return m.tryLock(ctx, goID, expiration)
	case <-m.PubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		return m.tryLock(ctx, goID, expiration)
	}
}

func (m *Mutex) lockInner(goID, expiration int64) (pTTL uint64, err error) {
	result, err := m.Client.Eval(context.TODO(), lockScript, []string{m.Name}, m.Uuid+":"+strconv.FormatInt(goID, 10), expiration).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return uint64(result.(int64)), nil
}

func (m *Mutex) Unlock() error {
	goID := base.GoID()
	return m.unlockInner(goID)
}

func (m *Mutex) unlockInner(goID int64) error {
	res, err := m.Client.Eval(context.TODO(), unlockScript, []string{m.Name, root.ChannelName(m.Name)}, m.Uuid+":"+strconv.FormatInt(goID, 10), 1).Int64()
	if err != nil {
		return err
	}
	if res == 0 {
		return errors.New("mismatch identification")
	}
	return err
}

func init() {
	path, _ := filepath.Abs(os.Args[1])
	index := strings.LastIndex(path, "/redisson")
	path = path[:index+9] + "/internal/mutex"

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
