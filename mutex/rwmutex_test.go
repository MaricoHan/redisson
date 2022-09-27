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
	rwMutex = NewRWMutex(&Root{
		Client: redis.NewClient(&redis.Options{Addr: ":6379"}),
		UUID:   "uuid",
	}, "rwMutexKey", []Option{
		WithExpireDuration(10 * time.Second),
		WithWaitTimeout(20 * time.Second),
	}...)
)

func TestRWMutex_lockInner(t *testing.T) {
	goID := util.GoID()
	pTTL, err := rwMutex.lockInner(goID, int64(rwMutex.expiration/time.Millisecond))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pTTL)
}

func TestRWMutex_unlockInner_ExpiredMutex(t *testing.T) {
	goID := util.GoID()

	// 测试：可以解锁过期的锁
	err := rwMutex.unlockInner(goID)
	if err != nil {
		t.Error(err)
		return
	}

	<-rwMutex.pubSub.Channel() // 测试：同样会发布解锁消息

	t.Log("unlock successfully")
}

func TestRWMutex_Unlock(t *testing.T) {
	err := rwMutex.Unlock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("unlock successfully")
}

func TestRWMutex_tryLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), rwMutex.waitTimeout)
	defer cancel()

	goID := util.GoID()
	err := rwMutex.tryLock(ctx, goID, int64(rwMutex.expiration/time.Millisecond))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("lock successfully")
}

func TestRWMutex_Lock(t *testing.T) {
	err := rwMutex.Lock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("lock successfully")
}

func TestRWMutex_rLockInner(t *testing.T) {
	goID := util.GoID()
	pTTL, err := rwMutex.rLockInner(goID, int64(rwMutex.expiration/time.Millisecond))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pTTL)
}

func TestRWMutex_tryRLock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), rwMutex.waitTimeout)
	defer cancel()

	goID := util.GoID()

	err := rwMutex.tryRLock(ctx, goID, int64(rwMutex.expiration/time.Millisecond))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("rLock successfully")
}

func TestRWMutex_RLock(t *testing.T) {
	err := rwMutex.RLock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("rLock successfully")
}

func TestRWMutex_Lock_RLock(t *testing.T) {
	TestRWMutex_Lock(t)

	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		defer group.Done()
		TestRWMutex_RLock(t) // “读锁”会阻塞到“写锁”释放，才会加锁成功，即 10s
	}()
	<-time.After(time.Second * 10)
	TestRWMutex_Unlock(t)

	group.Wait()
}

func TestMutex_Lock_Renewal(t *testing.T) {
	err := rwMutex.Lock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("lock successfully")

	// 测试：达到过期时间的 1/3，如果未主动释放锁，写锁的过期时间会被重置
	ticker := time.Tick(time.Second)
	for range ticker {
		fmt.Println(rwMutex.root.Client.PTTL(context.Background(), rwMutex.Name).Val())
	}

	time.After(2 * time.Minute)
}

func TestMutex_RLock_Renewal(t *testing.T) {
	err := rwMutex.RLock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("rLock successfully")

	// 测试：达到过期时间的 1/3，如果未主动释放锁，读锁的过期时间会被重置
	ticker := time.Tick(time.Second)
	for range ticker {
		fmt.Println(rwMutex.root.Client.PTTL(context.Background(), rwMutex.Name).Val())
	}

	time.After(2 * time.Minute)
}
