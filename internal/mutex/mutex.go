package mutex

import (
	"context"
	"github.com/MaricoHan/redisson/base"
	"github.com/MaricoHan/redisson/internal/root"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"strconv"
	"time"
)

var (
	lockScript   string
	unlockScript string
)

type Mutex struct {
	*root.Root
	Name    string
	TimeOut int64
}

func (m *Mutex) Lock() {
	//goID := base.GoID()

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
		return false, nil
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
			return false, nil
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

func (m *Mutex) Unlock() {

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
}
