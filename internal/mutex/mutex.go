package mutex

import (
	"context"
	"errors"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/MaricoHan/redisson/base"
	"github.com/MaricoHan/redisson/internal/root"
	"github.com/go-redis/redis/v8"
)

var (
	lockScript    string
	unlockScript  string
	renewalScript string
)

type Mutex struct {
	*root.Root
	Name    string
	TimeOut time.Duration
}

func (m *Mutex) Lock() error {
	_, err := m.tryLock()
	if err != nil {
		return err
	}
	// 加锁成功，开个协程，定时续锁
	go func() {
		ticker := time.Tick(time.Duration(m.TimeOut / 3))
		for range ticker {
			res, err := m.Client.Eval(context.TODO(), renewalScript, []string{m.Name}, m.TimeOut).Int64()
			if err != nil {
				return
			}
			if res == 0 {
				return
			}
		}

	}()

	return nil
}

func (m *Mutex) tryLock() (bool, error) {
	startTime := time.Now()

	goID := base.GoID()
	// 尝试加锁
	pTTL, err := m.tryAcquire(goID)
	if err != nil {
		return false, err
	}
	if pTTL == 0 {
		return true, nil
	}
	// 申请锁的耗时如果大于等于最大等待时间，则申请锁失败.
	if time.Now().After(startTime.Add(m.LockTimeout)) {
		return false, errors.New("timeout")
	}

	for {
		pTTL, err = m.tryAcquire(goID)
		if err != nil {
			return false, err
		}
		if pTTL == 0 {
			return true, nil
		}
		if time.Now().After(startTime.Add(m.LockTimeout)) {
			return false, errors.New("timeout")
		}
	}
}

func (m *Mutex) tryAcquire(goID int64) (pTTL uint64, err error) {
	result, err := m.Client.Eval(context.TODO(), lockScript, []string{m.Name}, m.Uuid+":"+strconv.FormatInt(goID, 10), m.TimeOut).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return uint64(result.(int64)), nil
}

func (m *Mutex) RLock() {

}

func (m *Mutex) Unlock() error {
	goID := base.GoID()
	if m.Uuid+":"+strconv.FormatInt(goID, 10) != m.Client.Get(context.TODO(), m.Name).String() {
		return errors.New("mismatch identification")
	}

	return m.unlockInner()
}

func (m *Mutex) unlockInner() error {
	_, err := m.Client.Eval(context.TODO(), unlockScript, []string{m.Name}).Result()
	if err == redis.Nil {
		return nil
	}
	return err
}

func init() {
	file, err := ioutil.ReadFile("./lock.lua")
	if err != nil {
		panic(err)
	}
	lockScript = string(file)

	file, err = ioutil.ReadFile("./unlock.lua")
	if err != nil {
		panic(err)
	}
	unlockScript = string(file)

	file, err = ioutil.ReadFile("./renewal.lua")
	if err != nil {
		panic(err)
	}
	renewalScript = string(file)
}
