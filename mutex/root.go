package mutex

import (
	"github.com/MaricoHan/redisson/pkg/utils/pubsub"
	"time"

	"github.com/go-redis/redis/v8"
)

type Root struct {
	Client *redis.Client
	UUID   string // 自定义用于区分不同客户端的唯一标识

	RedisChannelName string // redis 专用的 pubsub 频道名
}

type baseMutex struct {
	Name    string
	pubSub  *pubsub.PubSub
	release chan struct{}

	options *options
}

type options struct {
	expiration  time.Duration
	waitTimeout time.Duration
}

func (o *options) checkAndInit() {
	if o.waitTimeout <= 0 {
		o.waitTimeout = 30 * time.Second
	}
	if o.expiration <= 0 {
		o.expiration = 10 * time.Second
	}
}

type Option func(opts *options)

func WithExpireDuration(dur time.Duration) Option {
	return func(opt *options) {
		opt.expiration = dur
	}
}

func WithWaitTimeout(timeout time.Duration) Option {
	return func(opt *options) {
		opt.waitTimeout = timeout
	}
}
