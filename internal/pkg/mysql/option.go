package mysql

import (
	"io"
	"time"
)

type (
	config struct {
		maxIdleConns int
		maxOpenConns int
		maxLifetime  time.Duration
		debug        bool
		writer       io.Writer
	}

	Option func(cfg *config) error
)

func MaxIdleConnsOption(maxIdleConns int) Option {
	return func(cfg *config) error {
		if maxIdleConns == 0 {
			maxIdleConns = 10
		}
		cfg.maxIdleConns = maxIdleConns
		return nil
	}
}

func MaxOpenConnsOption(maxOpenConns int) Option {
	return func(cfg *config) error {
		if maxOpenConns == 0 {
			maxOpenConns = 10
		}
		cfg.maxOpenConns = maxOpenConns
		return nil
	}
}

func MaxLifetimeOption(maxLifetime string) Option {
	return func(cfg *config) error {
		maxLifetimeD, err := time.ParseDuration(maxLifetime)
		if err != nil {
			return err
		}

		if maxLifetimeD == 0 {
			maxLifetimeD = time.Hour
		}
		cfg.maxLifetime = maxLifetimeD
		return nil
	}
}

func DebugOption(debug bool) Option {
	return func(cfg *config) error {
		cfg.debug = debug
		return nil
	}
}

func WriteOption(writer io.Writer) Option {
	return func(cfg *config) error {
		cfg.writer = writer
		return nil
	}
}
