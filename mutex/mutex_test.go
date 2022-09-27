package mutex

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson/pkg/util"
)

var (
	mutex = NewMutex(&Root{
		Client: redis.NewClient(&redis.Options{Addr: ":6379"}),
		UUID:   "uuid",
	}, "mutexKey", []Option{
		WithExpireDuration(10 * time.Second),
		WithWaitTimeout(20 * time.Second),
	}...)
)

func TestMutex_lockInner(t *testing.T) {
	acquire, err := mutex.lockInner(util.GoID(), int64(mutex.Expiration/time.Millisecond))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(acquire)
}

func TestMutex_tryLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), mutex.WaitTimeout)
	defer cancel()
	goID := util.GoID()

	err := mutex.tryLock(ctx, goID, int64(mutex.Expiration/time.Millisecond))
	t.Log(err)
}

func TestMutex_unlockInner_ExpiredMutex(t *testing.T) {
	goID := util.GoID()

	// 测试：可以解锁过期的锁
	err := mutex.unlockInner(goID)
	if err != nil {
		t.Error(err)
		return
	}

	<-mutex.PubSub.Channel() // 测试：同样会发布解锁消息

	t.Log("unlock successfully")
}

// TestMutex_unlockInner
// @Description: 测试：只可以解自己加的锁
// @param t
func TestMutex_unlockInner(t *testing.T) {
	goID := util.GoID()

	_, err := mutex.lockInner(goID, int64(mutex.Expiration/time.Millisecond))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("lock successfully")

	// 测试：其他协程无法解锁
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer func() {
			waitGroup.Done()
		}()
		err = mutex.unlockInner(util.GoID())
		if err != nil {
			t.Error(err) // mismatch identification
			return
		}
		t.Log("unlock successfully")
	}()
	waitGroup.Wait()

	// 测试：加锁的协程可以顺利解锁
	err = mutex.unlockInner(goID)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("unlock successfully")
}

func TestMutex_Unlock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), mutex.WaitTimeout)
	defer cancel()
	goID := util.GoID()
	// 第一次上锁
	err := mutex.tryLock(ctx, goID, int64(mutex.Expiration/time.Millisecond))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("lock successfully")

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), mutex.WaitTimeout)
		defer func() {
			cancel()
			waitGroup.Done()
		}()
		// 不解锁，第二次上锁，会阻塞 10s，然后加锁成功
		t.Log("try lock ...")
		err = mutex.tryLock(ctx, goID, int64(mutex.Expiration/time.Millisecond))
		cancel()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("lock successfully")
	}()

	// 10s 后解锁
	timer := time.NewTimer(10 * time.Second)
	<-timer.C
	err = mutex.Unlock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("unlock successfully")

	waitGroup.Wait()
}

func TestMutex_Renewal(t *testing.T) {
	err := mutex.Lock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("lock successfully")

	// 测试：达到过期时间的 1/3，如果未主动释放锁，锁的过期时间会被重置
	ticker := time.Tick(time.Second)
	for range ticker {
		fmt.Println(mutex.root.Client.PTTL(context.Background(), mutex.Name).Val())
	}

	time.After(2 * time.Minute)
}
