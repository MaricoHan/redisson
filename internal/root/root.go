package root

import (
	"github.com/MaricoHan/redisson"
	"github.com/go-redis/redis/v8"
)

type Root struct {
	Client *redis.Client
	Uuid   string
	*redisson.Options
}
