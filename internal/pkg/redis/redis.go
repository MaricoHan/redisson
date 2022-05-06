package redis

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"time"
)

type RedisClient struct {
	client *redis.Client
	log    *log.Logger
}

func NewRedisClient(addr, password string, db int64, log *log.Logger) *RedisClient {

	client := redis.NewClient(&redis.Options{
		Addr:     addr,     // use default Addr
		Password: password, // no password set
		DB:       int(db),  // use default DB
	})
	return &RedisClient{
		client: client,
		log:    log,
	}

}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(context.Background(), keyPrefix(key), value, expiration).Err()
}

func (r *RedisClient) Ping() bool {
	err := r.client.Ping(context.Background()).Err()
	if err != nil {
		return false
	}
	return true
}

func (r *RedisClient) Close() {
	r.client.Close()
	return
}

func (r *RedisClient) Has(key string) (bool, error) {
	result, err := r.client.Exists(context.Background(), keyPrefix(key)).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (r *RedisClient) Delete(key string) error {
	return r.client.Del(context.Background(), keyPrefix(key)).Err()
}

func keyPrefix(key string) string {
	return fmt.Sprintf("%s:%s", constant.RedisPrefix, key)
}
