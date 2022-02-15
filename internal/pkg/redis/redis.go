package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bsm/redislock"
	redis "github.com/go-redis/redis/v8"
)

const prefix = "nftp"

var (
	rdb    *redis.Client
	locker *redislock.Client

	ErrNotObtained = redislock.ErrNotObtained
)

func RedisPing() bool {
	conn := rdb.Conn(context.Background())
	if conn == nil {
		return false
	}
	return true
}

// Connect connect tht redis server
func Connect(addr, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	// Create a new lock client.
	locker = redislock.New(rdb)
}

// Obtain try to obtain lock.
func Obtain(ctx context.Context, key string, ttl time.Duration) (*redislock.Lock, error) {
	return locker.Obtain(ctx, keyPrefix(key), ttl, nil)
}

// Set save a (key,value) to redis
func Set(key string, value interface{}, expiration time.Duration) error {
	return rdb.Set(context.Background(), keyPrefix(key), value, expiration).Err()
}

// Has return a bool value
func Has(key string) (bool, error) {
	result, err := rdb.Exists(context.Background(), keyPrefix(key)).Result()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

// GetString return a string value from redis by specisal key
func GetString(key string) (string, error) {
	return rdb.Get(context.Background(), keyPrefix(key)).Result()
}

// GetInt64 return a int64 value from redis by specisal key
func GetInt64(key string) (int64, error) {
	return rdb.Get(context.Background(), keyPrefix(key)).Int64()
}

// SetObject save a object value to redis by specisal key
func SetObject(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(context.Background(), keyPrefix(key), data, expiration).Err()
}

// GetObject return a object value from redis by specisal key
func GetObject(key string, value interface{}) error {
	data, err := rdb.Get(context.Background(), keyPrefix(key)).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

// HSet save a simple value to redis hash table by specisal key
func HSet(key string, values ...interface{}) error {
	return rdb.HSet(context.Background(), keyPrefix(key), values...).Err()
}

// HExists returns whether the field exists
func HExists(key string, field string) (bool, error) {
	if key == "" || field == "" {
		return false, errors.New("参数为空")
	}
	has, err := rdb.HExists(context.Background(), keyPrefix(key), field).Result()
	if err != nil {
		return false, errors.New("异常：" + err.Error())
	}
	return has, nil
}

// Delete remove a value from redis by specisal key
func Delete(key string) error {
	return rdb.Del(context.Background(), keyPrefix(key)).Err()
}

// Has return a bool value
func Close() {
	p := fmt.Sprintf("%s*", prefix)

	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, p, 0).Iterator()
	for iter.Next(ctx) {
		err := rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			panic(err)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	_ = rdb.Close()
}

// Publish publish message to `channel`
func Publish(channel string, payload string) error {
	return rdb.Publish(context.Background(), channel, payload).Err()
}

func Pipeline() redis.Pipeliner {
	return rdb.Pipeline()
}

// Subscribe subscribe message from `channel`
func Subscribe(channel string,
	success func(message string) error,
	fail func(err error),
) {
	go func() {
		pubsub := rdb.Subscribe(context.Background(), channel)
		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				fail(err)
				return
			}

			if err = success(msg.Payload); err != nil {
				fail(err)
			}
		}
	}()
}

func keyPrefix(key string) string {
	return fmt.Sprintf("%s:%s", prefix, key)
}
