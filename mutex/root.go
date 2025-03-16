package mutex

import (
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson/pkg/loggers"
	"github.com/MaricoHan/redisson/pkg/utils/pubsub"
)

// Root 是所有锁的根结构，包含共享资源
type Root struct {
	Client *redis.Client
	UUID   string // 自定义用于区分不同客户端的唯一标识

	RedisChannelName string           // redis 专用的 pubsub 频道名
	Logger           loggers.Advanced // 日志接口
}

// baseMutex 是所有锁类型的基础结构
type baseMutex struct {
	Name    string
	pubSub  *pubsub.PubSub
	release chan struct{}

	options *options
}

// options 定义锁的配置选项
type options struct {
	expiration  time.Duration // 锁的过期时间
	waitTimeout time.Duration // 获取锁的最大等待时间
}

// checkAndInit 检查并初始化选项的默认值
func (o *options) checkAndInit() {
	if o.waitTimeout <= 0 {
		o.waitTimeout = 30 * time.Second
	}
	if o.expiration <= 0 {
		o.expiration = 10 * time.Second
	}
}

// Option 是配置锁选项的函数类型
type Option func(opts *options)

// WithExpireDuration 设置锁的过期时间
func WithExpireDuration(dur time.Duration) Option {
	return func(opt *options) {
		opt.expiration = dur
	}
}

// WithWaitTimeout 设置获取锁的最大等待时间
func WithWaitTimeout(timeout time.Duration) Option {
	return func(opt *options) {
		opt.waitTimeout = timeout
	}
}
