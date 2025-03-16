package redisson

import (
	"context"
	"runtime"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"

	"github.com/MaricoHan/redisson/mutex"
	"github.com/MaricoHan/redisson/pkg/loggers"
	"github.com/MaricoHan/redisson/pkg/utils"
	"github.com/MaricoHan/redisson/pkg/utils/pubsub"
)

type Redisson struct {
	root *mutex.Root
}

type Config struct {
	Logger loggers.Advanced
}

func DefaultConfig() *Config {
	return &Config{
		Logger: nil,
	}
}

func (c *Config) CheckAndInit() {
	if c.Logger == nil {
		c.Logger = loggers.Logger()
	}
}

func New(ctx context.Context, client *redis.Client) *Redisson {
	return NewWithConfig(ctx, client, DefaultConfig())
}

func NewWithConfig(ctx context.Context, client *redis.Client, config *Config) *Redisson {
	config.CheckAndInit()

	redisson := &Redisson{
		root: &mutex.Root{
			Client:           client,
			UUID:             ksuid.New().String(),
			RedisChannelName: utils.ChannelName("redisson_pubsub"),
			Logger:           config.Logger,
		},
	}

	config.Logger.Debugf("初始化 Redisson 实例，UUID: %s", redisson.root.UUID)

	// 一个实例只建立一个 pubsub 连接
	// 额外开协程监听 redis 消息，转发给实例内部基于内存实现的 pubsub，再分配给对应的 subscriber。
	pubSub := client.Subscribe(ctx, redisson.root.RedisChannelName)

	gCtx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		defer func() { _ = pubSub.Close() }()

		var ss []string
		for {
			select {
			case <-gCtx.Done():
				config.Logger.Debug("Redisson pubsub 监听协程退出")
				return
			case msg := <-pubSub.Channel():
				ss = strings.Split(msg.Payload, ":")
				config.Logger.Debugf("收到 Redis 消息：%s, 动作：%s", ss[0], ss[1])
				pubsub.Publish(ss[0], ss[1]) // 0:锁名 1:动作(unlock)
			}
		}
	}()
	wg.Wait() // 等待协程启动成功

	// 释放资源
	runtime.SetFinalizer(redisson, func(_ *Redisson) {
		config.Logger.Debug("释放 Redisson 资源")
		cancel()
	})

	return redisson
}

func (r Redisson) NewMutex(name string, options ...mutex.Option) *mutex.Mutex {
	return mutex.NewMutex(r.root, name, options...)
}

func (r Redisson) NewRWMutex(name string, options ...mutex.Option) *mutex.RWMutex {
	return mutex.NewRWMutex(r.root, name, options...)
}
