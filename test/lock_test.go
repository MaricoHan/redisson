package test

import (
	"github.com/MaricoHan/redisson"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	redisson := redisson.New(client)

	mutex := redisson.NewRWMutex("han")
	mutex.Lock()
	mutex.Unlock()
}
