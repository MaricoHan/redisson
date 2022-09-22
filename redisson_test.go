package redisson_test

import (
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson"
	"github.com/MaricoHan/redisson/internal/root"
)

func TestMutex(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	options := []redisson.Option{
		redisson.WithExpireDuration(30 * time.Millisecond),
	}

	redisson := redisson.New(client, &root.Options{
		LockTimeout: 10 * time.Second,
	})

	mutex := redisson.NewMutex("redisson_mutex", options...)

	err := mutex.Lock()
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
		mutex := redisson.NewMutex("redisson_mutex")
		err = mutex.Unlock()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("unlock successfully")
	}()
	waitGroup.Wait()

	// 测试：加锁的协程可以顺利解锁
	err = mutex.Unlock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("unlock successfully")
}
