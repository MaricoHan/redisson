package redisson

import "time"

type Options struct {
	LockTimeout time.Duration
}

func (o *Options) init() {
	if o.LockTimeout <= 0 {
		o.LockTimeout = 30 * time.Second
	}
}
