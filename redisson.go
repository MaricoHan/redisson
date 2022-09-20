package redisson

import (
	"github.com/MaricoHan/redisson/internal/mutex"
	"github.com/MaricoHan/redisson/internal/root"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"
)

type Redisson struct {
	root *root.Root
}

func New(client *redis.Client, options *Options) *Redisson {
	options.init()

	return &Redisson{
		root: &root.Root{
			Client:  client,
			Options: options,
			Uuid:    ksuid.New().String(),
		},
	}
}

func (r *Redisson) NewMutex(name string, options ...Option) *mutex.Mutex {

	m := &mutex.Mutex{Root: r.root, Name: name}

	for i := range options {
		options[i].Apply(m)
	}

	return m
}

func (r *Redisson) NewRWMutex() {

}

type Option interface {
	Apply(mutex *mutex.Mutex)
}

type OptionFunc func(mutex *mutex.Mutex)

func (f OptionFunc) Apply(mutex *mutex.Mutex) {
	f(mutex)
}

func WithExpireDuration(dur int64) Option {
	return OptionFunc(func(mutex *mutex.Mutex) {
		mutex.TimeOut = dur
	})
}
