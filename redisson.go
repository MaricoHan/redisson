package redisson

import (
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"

	"github.com/MaricoHan/redisson/mutex"
)

type Redisson struct {
	root *mutex.Root
}

func New(client *redis.Client) *Redisson {
	return &Redisson{
		root: &mutex.Root{
			Client: client,
			UUID:   ksuid.New().String(),
		},
	}
}

func (r Redisson) NewMutex(name string, options ...mutex.Option) *mutex.Mutex {
	return mutex.NewMutex(r.root, name, options...)
}

func (r Redisson) NewRWMutex(name string, options ...mutex.Option) *mutex.RWMutex {
	return mutex.NewRWMutex(r.root, name, options...)
}
