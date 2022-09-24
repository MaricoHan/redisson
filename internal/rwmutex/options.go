package rwmutex

import "time"

type Option interface {
	Apply(mutex *RWMutex)
}

type OptionFunc func(mutex *RWMutex)

func (f OptionFunc) Apply(mutex *RWMutex) {
	f(mutex)
}

func WithExpireDuration(dur time.Duration) Option {
	return OptionFunc(func(mutex *RWMutex) {
		mutex.Expiration = dur
	})
}

func WithWaitTimeout(timeout time.Duration) Option {
	return OptionFunc(func(mutex *RWMutex) {
		mutex.WaitTimeout = timeout
	})
}
