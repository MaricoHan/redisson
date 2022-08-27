package redisson

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type Redisson struct {
	client *redis.Client
}

func New(client *redis.Client) *Redisson {
	return &Redisson{
		client: client,
	}
}

func (r *Redisson) NewRWMutex(name string, options ...Option) *Mutex {

	m := &Mutex{client: r.client, Name: name}

	for i := range options {
		options[i].Apply(m)
	}

	return m
}

type Option interface {
	Apply(mutex *Mutex)
}

type OptionFunc func(mutex *Mutex)

func (f OptionFunc) Apply(mutex *Mutex) {
	f(mutex)
}

func WithExpireDuration(dur time.Duration) Option {
	return OptionFunc(func(mutex *Mutex) {
		mutex.TimeOut = dur
	})
}
