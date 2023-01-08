package mutex

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Root struct {
	Client *redis.Client
	UUID   string // 自定义用于区分不同客户端的唯一标识
}

type baseMutex struct {
	Name        string
	expiration  time.Duration
	waitTimeout time.Duration
	pubSub      *redis.PubSub
}

func (b *baseMutex) checkAndInit() {
	if b.waitTimeout <= 0 {
		b.waitTimeout = 30 * time.Second
	}
	if b.expiration <= 0 {
		b.expiration = 10 * time.Second
	}
}

type Option interface {
	Apply(mutex *baseMutex)
}

type OptionFunc func(mutex *baseMutex)

func (f OptionFunc) Apply(mutex *baseMutex) {
	f(mutex)
}

func WithExpireDuration(dur time.Duration) Option {
	return OptionFunc(func(mutex *baseMutex) {
		mutex.expiration = dur
	})
}

func WithWaitTimeout(timeout time.Duration) Option {
	return OptionFunc(func(mutex *baseMutex) {
		mutex.waitTimeout = timeout
	})
}
