package redisson

import (
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"

	"github.com/MaricoHan/redisson/internal/mutex"
	"github.com/MaricoHan/redisson/internal/root"
	"github.com/MaricoHan/redisson/internal/rwmutex"
)

type Redisson struct {
	root *root.Root
}

func New(client *redis.Client) *Redisson {
	return &Redisson{
		root: &root.Root{
			Client: client,
			UUID:   ksuid.New().String(),
		},
	}
}

func (r Redisson) NewMutex(name string, options ...mutex.Option) *mutex.Mutex {
	return mutex.NewMutex(r.root, name, options...)
}

func (r Redisson) NewRWMutex(name string, options ...rwmutex.Option) *rwmutex.RWMutex {
	return rwmutex.NewRWMutex(r.root, name, options...)
}
