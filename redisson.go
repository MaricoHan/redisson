package redisson

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"

	"github.com/MaricoHan/redisson/internal/mutex"
	"github.com/MaricoHan/redisson/internal/root"
)

type Redisson struct {
	root *root.Root
}

func New(client *redis.Client, options *root.Options) *Redisson {
	options.Init()

	return &Redisson{
		root: &root.Root{
			Client:  client,
			Options: options,
			Uuid:    ksuid.New().String(),
		},
	}
}

func (r *Redisson) NewMutex(name string, options ...Option) *mutex.Mutex {
	m := &mutex.Mutex{
		Root:   r.root,
		Name:   name,
		PubSub: r.root.Client.Subscribe(context.Background(), root.ChannelName(name)),
	}

	for i := range options {
		options[i].Apply(m)
	}

	return m.Init()
}

type Option interface {
	Apply(mutex *mutex.Mutex)
}

type OptionFunc func(mutex *mutex.Mutex)

func (f OptionFunc) Apply(mutex *mutex.Mutex) {
	f(mutex)
}

func WithExpireDuration(dur time.Duration) Option {
	return OptionFunc(func(mutex *mutex.Mutex) {
		mutex.Expiration = dur
	})
}
