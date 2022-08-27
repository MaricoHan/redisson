package root

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type Root struct {
	Client *redis.Client
	Uuid   string
	*Options
}

func ChannelName(name string) string {
	return "redisson_lock__channel" + ":{" + name + "}"
}

type Options struct {
	LockTimeout time.Duration
}

func (o *Options) Init() {
	if o.LockTimeout <= 0 {
		o.LockTimeout = 30 * time.Second
	}
}
