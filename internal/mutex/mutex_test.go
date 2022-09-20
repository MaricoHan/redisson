package mutex

import (
	"fmt"
	"github.com/MaricoHan/redisson/base"
	"github.com/MaricoHan/redisson/internal/root"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"testing"
)

func TestReadFile(t *testing.T) {
	file, err := ioutil.ReadFile("./lock.lua")
	if err != nil {
		return
	}
	fmt.Println(string(file))
}

var mutex = Mutex{
	Root: &root.Root{
		Client: redis.NewClient(&redis.Options{Addr: ":6379"}),
		Uuid:   "uuid",
	},
	Name:    "mutexKey",
	TimeOut: 1000000,
}

func TestMutex_tryAcquire(t *testing.T) {
	acquire, err := mutex.tryAcquire(base.GoID())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(acquire)
}
