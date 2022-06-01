package redis

import (
	"context"
	"encoding/json"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"strings"
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

// SetObject save a object value to redis by special key
func (r *RedisClient) SetObject(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(context.Background(), keyPrefix(key), data, expiration).Err()
}

// GetObject return a object value from redis by special key
func (r *RedisClient) GetObject(key string, value interface{}) error {
	data, err := r.client.Get(context.Background(), keyPrefix(key)).Bytes()
	if err != nil && !strings.Contains(err.Error(), "redis: nil") {
		return err
	}
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, value)
}
