package redisson

import (
	"context"
	"runtime"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/segmentio/ksuid"

	"github.com/MaricoHan/redisson/mutex"
	"github.com/MaricoHan/redisson/pkg/utils"
	"github.com/MaricoHan/redisson/pkg/utils/pubsub"
)

type Redisson struct {
	root *mutex.Root
}

func New(ctx context.Context, client *redis.Client) *Redisson {
	redisson := &Redisson{
		root: &mutex.Root{
			Client:           client,
			UUID:             ksuid.New().String(),
			RedisChannelName: utils.ChannelName("redisson_pubsub"),
		},
	}

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
				return
			case msg := <-pubSub.Channel():
				ss = strings.Split(msg.Payload, ":")
				pubsub.Publish(ss[0], ss[1]) // 0:锁名 1:动作(unlock)
			}
		}
	}()
	wg.Wait() // 等待协程启动成功

	// 释放资源
	runtime.SetFinalizer(redisson, func(_ *Redisson) { cancel() })

	return redisson
}

func (r Redisson) NewMutex(name string, options ...mutex.Option) *mutex.Mutex {
	return mutex.NewMutex(r.root, name, options...)
}

func (r Redisson) NewRWMutex(name string, options ...mutex.Option) *mutex.RWMutex {
	return mutex.NewRWMutex(r.root, name, options...)
}
