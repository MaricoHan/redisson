package redisson

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Mutex struct {
	client  *redis.Client
	Name    string
	TimeOut time.Duration
}

func (m *Mutex) Lock() {

}

func (m *Mutex) tryLock() {

}

func (m *Mutex) tryAcquire() {
	result, err := m.client.Eval(context.TODO(), `if (redis.call('exists',KEYS[1]) == 0) then
    redis.call('hincrby',KEYS[1],ARGV[2],1);
    redis.call('pexpire',KEYS[1],ARGV[1]);
    return nil;
end
if (redis.call('hexist',KEYS[1],ARGV[2])==1) then
    redis.call('hincrby',KEYS[1],ARGV[2],1);
    redis.call('pexpire',KEYS[1],ARGV[1]);
    return nil;
end
return redis.call('pttl',KEYS[1]);`, []string{m.Name}, m.TimeOut).Result()
	if err != nil {
		return
	}

}

func (m *Mutex) RLock() {

}

func (m *Mutex) Unlock() {

}
