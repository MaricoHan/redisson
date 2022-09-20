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

type Options struct {
	LockTimeout time.Duration
}

func (o *Options) init() {
	if o.LockTimeout <= 0 {
		o.LockTimeout = 30 * time.Second
	}
}
