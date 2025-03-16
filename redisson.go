package redisson

import (
	"context"
	"runtime"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

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
			UUID:             uuid.New().String(),
			RedisChannelName: utils.ChannelName("redisson_pubsub"),
			Logger:           config.Logger,
		},
	}

	config.Logger.Infof("初始化 Redisson 实例，UUID: %s, Redis通道: %s", redisson.root.UUID, redisson.root.RedisChannelName)

	// 一个实例只建立一个 pubsub 连接
	// 额外开协程监听 redis 消息，转发给实例内部基于内存实现的 pubsub，再分配给对应的 subscriber。
	pubSub := client.Subscribe(ctx, redisson.root.RedisChannelName)
	config.Logger.Debugf("订阅 Redis 通道: %s", redisson.root.RedisChannelName)

	gCtx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		defer func() {
			err := pubSub.Close()
			if err != nil {
				config.Logger.Errorf("关闭 Redis 订阅连接失败: %v", err)
			} else {
				config.Logger.Debug("关闭 Redis 订阅连接成功")
			}
		}()

		var ss []string
		config.Logger.Info("启动 Redis 消息监听协程")
		for {
			select {
			case <-gCtx.Done():
				config.Logger.Info("Redisson pubsub 监听协程收到退出信号")
				return
			case msg := <-pubSub.Channel():
				ss = strings.Split(msg.Payload, ":")
				config.Logger.Debugf("收到 Redis 消息: %s, 动作: %s, 通道: %s", ss[0], ss[1], msg.Channel)
				pubsub.Publish(ss[0], ss[1]) // 0:锁名 1:动作(unlock)
			}
		}
	}()
	wg.Wait() // 等待协程启动成功
	config.Logger.Debug("Redis 消息监听协程启动成功")

	// 释放资源
	runtime.SetFinalizer(redisson, func(_ *Redisson) {
		config.Logger.Info("释放 Redisson 资源，UUID: " + redisson.root.UUID)
		cancel()
	})

	return redisson
}

func (r Redisson) NewMutex(name string, options ...mutex.Option) *mutex.Mutex {
	r.root.Logger.Debugf("创建互斥锁: %s", name)
	return mutex.NewMutex(r.root, name, options...)
}

func (r Redisson) NewRWMutex(name string, options ...mutex.Option) *mutex.RWMutex {
	r.root.Logger.Debugf("创建读写锁: %s", name)
	return mutex.NewRWMutex(r.root, name, options...)
}
