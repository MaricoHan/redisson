package mutex

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Root struct {
	Client *redis.Client
	UUID   string
}

type baseMutex struct {
	Name        string
	Expiration  time.Duration
	WaitTimeout time.Duration
	*redis.PubSub
}

func (b baseMutex) CheckAndInit() {
	if b.WaitTimeout <= 0 {
		b.WaitTimeout = 30 * time.Second
	}
	if b.Expiration <= 0 {
		b.Expiration = 10 * time.Second
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
		mutex.Expiration = dur
	})
}

func WithWaitTimeout(timeout time.Duration) Option {
	return OptionFunc(func(mutex *baseMutex) {
		mutex.WaitTimeout = timeout
	})
}
