package redisson_test

import (
	"sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson"
	"github.com/MaricoHan/redisson/mutex"
)

func TestMutex(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	redissonClient := redisson.New(client)

	options := []mutex.Option{
		mutex.WithExpireDuration(30 * time.Millisecond),
	}
	mutex1 := redissonClient.NewMutex("redisson_mutex", options...)

	err := mutex1.Lock()
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
		var mutex2 = redissonClient.NewMutex("redisson_mutex")
		err = mutex2.Unlock()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("unlock successfully")
	}()
	waitGroup.Wait()

	// 测试：加锁的协程可以顺利解锁
	err = mutex1.Unlock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("unlock successfully")
}

func TestRWMutex(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	redissonClient := redisson.New(client)

	options := []mutex.Option{
		mutex.WithExpireDuration(30 * time.Millisecond),
	}
	mutex1 := redissonClient.NewRWMutex("redisson_mutex", options...)

	err := mutex1.Lock()
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
		var mutex2 = redissonClient.NewMutex("redisson_mutex")
		err = mutex2.Unlock()
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("unlock successfully")
	}()
	waitGroup.Wait()

	// 测试：加锁的协程可以顺利解锁
	err = mutex1.Unlock()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("unlock successfully")
}
