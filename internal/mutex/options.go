package mutex

import "time"

type Option interface {
	Apply(mutex *Mutex)
}

type OptionFunc func(mutex *Mutex)

func (f OptionFunc) Apply(mutex *Mutex) {
	f(mutex)
}

func WithExpireDuration(dur time.Duration) Option {
	return OptionFunc(func(mutex *Mutex) {
		mutex.Expiration = dur
	})
}

func WithWaitTimeout(timeout time.Duration) Option {
	return OptionFunc(func(mutex *Mutex) {
		mutex.WaitTimeout = timeout
	})
}
