package rwmutex

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MaricoHan/redisson/internal/root"
	"github.com/MaricoHan/redisson/pkg/base"
	"github.com/MaricoHan/redisson/pkg/types"
)

var (
	lockScript    = ""
	rLockScript   = ""
	renewalScript = ""
	unlockScript  = ""
)

type RWMutex struct {
	*root.Root
	root.BaseMutex
}

func NewRWMutex(r *root.Root, name string, options ...Option) *RWMutex {
	baseMutex := root.BaseMutex{
		Name:   name,
		PubSub: r.Client.Subscribe(context.Background(), base.ChannelName(name)),
	}

	m := &RWMutex{
		Root:      r,
		BaseMutex: baseMutex,
	}

	for i := range options {
		options[i].Apply(m)
	}

	return m.CheckAndInit()
}

func (r *RWMutex) CheckAndInit() *RWMutex {
	r.BaseMutex.CheckAndInit()

	return r
}

func (r RWMutex) Lock() error {
	// 单位：ms
	expiration := int64(r.Expiration / time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), r.WaitTimeout)
	defer cancel()
	goID := base.GoID()
	err := r.tryLock(ctx, goID, expiration)
	if err != nil {
		return err
	}
	// 加锁成功，开个协程，定时续锁
	go func() {
		ticker := time.NewTicker(r.Expiration / 3).C
		for range ticker {
			res, err := r.Client.Eval(context.TODO(), renewalScript, []string{r.Name}, expiration).Int64()
			if err != nil || res == 0 {
				return
			}
		}
	}()

	return nil

}

func (r RWMutex) tryLock(ctx context.Context, goID, expiration int64) error {
	pTTL, err := r.lockInner(goID, expiration)
	if err != nil {
		return err
	}

	if pTTL == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return types.ErrWaitTimeout
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对“redis 中存在未维护的锁”，即当锁自然过期后，并不会发布通知的锁
		return r.tryLock(ctx, goID, expiration)
	case <-r.PubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		return r.tryLock(ctx, goID, expiration)
	}

}

func (r RWMutex) lockInner(goID, expiration int64) (int64, error) {
	pTTL, err := r.Client.Eval(context.Background(), lockScript, []string{r.Name}, r.UUID+":"+strconv.FormatInt(goID, 10), expiration).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return pTTL.(int64), nil
}

func (r RWMutex) RLock() error {
	// 单位：ms
	expiration := int64(r.Expiration / time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), r.WaitTimeout)
	defer cancel()
	goID := base.GoID()
	err := r.tryRLock(ctx, goID, expiration)
	if err != nil {
		return err
	}
	// 加锁成功，开个协程，定时续锁
	go func() {
		ticker := time.NewTicker(r.Expiration / 3).C
		for range ticker {
			res, err := r.Client.Eval(context.TODO(), renewalScript, []string{r.Name}, expiration).Int64()
			if err != nil || res == 0 {
				return
			}
		}
	}()

	return nil

}

func (r RWMutex) tryRLock(ctx context.Context, goID, expiration int64) error {
	pTTL, err := r.rLockInner(goID, expiration)
	if err != nil {
		return err
	}

	if pTTL == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return types.ErrWaitTimeout
	case <-time.After(time.Duration(pTTL) * time.Millisecond):
		// 针对“redis 中存在未维护的锁”，即当锁自然过期后，并不会发布通知的锁
		return r.tryRLock(ctx, goID, expiration)
	case <-r.PubSub.Channel():
		// 收到解锁通知，则尝试抢锁
		return r.tryRLock(ctx, goID, expiration)
	}

}

func (r RWMutex) rLockInner(goID, expiration int64) (int64, error) {
	pTTL, err := r.Client.Eval(context.Background(), rLockScript, []string{r.Name}, r.UUID+":"+strconv.FormatInt(goID, 10), expiration).Result()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return pTTL.(int64), nil
}
func (r *RWMutex) Unlock() error {
	goID := base.GoID()
	return r.unlockInner(goID)
}
func (r *RWMutex) unlockInner(goID int64) error {
	res, err := r.Client.Eval(context.TODO(), unlockScript, []string{r.Name, base.ChannelName(r.Name)}, r.UUID+":"+strconv.FormatInt(goID, 10), 1).Int64()
	if err != nil {
		return err
	}
	if res == 0 {
		return types.ErrMismatch
	}

	return nil
}

func init() {
	path, _ := filepath.Abs(os.Args[1])
	index := strings.LastIndex(path, "/redisson")
	path = path[:index+9] + "/internal/rwmutex/lua/"

	file, err := ioutil.ReadFile(path + "lock.lua")
	if err != nil {
		panic(err)
	}
	lockScript = string(file)

	file, err = ioutil.ReadFile(path + "rlock.lua")
	if err != nil {
		panic(err)
	}
	rLockScript = string(file)

	file, err = ioutil.ReadFile(path + "renewal.lua")
	if err != nil {
		panic(err)
	}
	renewalScript = string(file)

	file, err = ioutil.ReadFile(path + "unlock.lua")
	if err != nil {
		panic(err)
	}
	unlockScript = string(file)
}
