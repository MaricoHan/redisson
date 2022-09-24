package root

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Root struct {
	Client *redis.Client
	Uuid   string
}

type BaseMutex struct {
	Name        string
	Expiration  time.Duration
	WaitTimeout time.Duration
	*redis.PubSub
}

func (b *BaseMutex) CheckAndInit() {
	if b.WaitTimeout <= 0 {
		b.WaitTimeout = 30 * time.Second
	}
	if b.Expiration <= 0 {
		b.Expiration = 10 * time.Second
	}
}
