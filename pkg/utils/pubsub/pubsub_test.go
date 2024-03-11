package pubsub

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	// 启动 10 个订阅者
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			pubSub := Subscribe("channel_" + fmt.Sprintf("%d", i))

			wg.Done()
			for {
				select {
				case msg := <-pubSub.Channel():
					fmt.Println("channel_"+fmt.Sprintf("%d", i)+" receive: ", msg)
				}
			}
		}(i)
	}
	wg.Wait()

	// 发布 n 条消息
	for i := 0; i < 22; i++ {
		Publish("channel_"+fmt.Sprintf("%d", i%10), "hello "+fmt.Sprintf("%d", i))
	}

	time.Sleep(time.Second)
}
