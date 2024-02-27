package pubsub

import (
	"sync"
	"time"
)

var (
	channels = make(map[string][]*PubSub)

	mu = sync.Mutex{}
)

type PubSub struct {
	channelName        string
	msgChan            chan string
	msgChanSize        int
	msgChanSendTimeout time.Duration

	chOnce sync.Once
}

func (p *PubSub) Channel() <-chan string {
	p.chOnce.Do(func() {
		// 设置默认值，后期可改为 option 模式作为入参
		p.msgChanSize = 100
		p.msgChanSendTimeout = time.Second

		p.msgChan = make(chan string, p.msgChanSize)
	})

	return p.msgChan
}

func (p *PubSub) Close() {
	mu.Lock()
	defer mu.Unlock()

	if p.msgChan != nil {
		close(p.msgChan)
	}

	subs, exists := channels[p.channelName]
	if !exists {
		return
	}

	// 查找并删除对应的订阅者
	for i := range subs {
		if subs[i] == p {
			subs = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	// 如果没有订阅者了，则删除对应的频道
	if len(subs) == 0 {
		delete(channels, p.channelName)
	}
}

func Subscribe(channelName string) *PubSub {
	mu.Lock()
	defer mu.Unlock()

	pb := &PubSub{
		channelName: channelName,
		chOnce:      sync.Once{},
	}

	channels[channelName] = append(channels[channelName], pb)

	return pb
}

func Publish(channelName string, msg string) {
	mu.Lock()
	defer mu.Unlock()

	subscribers, exists := channels[channelName]
	if !exists {
		return
	}

	// 非阻塞发布
	for _, sub := range subscribers {
		select {
		case sub.msgChan <- msg:
		case <-time.After(sub.msgChanSendTimeout):
		}
	}
}
