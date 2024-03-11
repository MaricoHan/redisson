# redisson

[![Go](https://github.com/MaricoHan/redisson/actions/workflows/go.yml/badge.svg)](https://github.com/MaricoHan/redisson/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/MaricoHan/redisson)](https://goreportcard.com/report/github.com/MaricoHan/redisson)
[![License：MIT](https://img.shields.io/github/license/MaricoHan/redisson)](https://github.com/MaricoHan/redisson/blob/master/LICENSE)

基于 redis 实现：分布式“互斥锁”和“读写锁”。

目前不支持可重入锁、红锁。

# 功能

## 互斥锁

* 任一时刻，对于一把锁，只可以有一个持有者。

## 读写锁

* 读锁与读锁可以共存，写锁与读锁/写锁不可以共存。

> 无论是互斥锁还是读写锁，都只可以由加锁的协程解锁，其他协程无法解锁。
> 加锁成功以后会开启一个协程定时续锁，直到客户端解锁。

# 使用

## 获取依赖

```shell
go get github.com/MaricoHan/redisson
```

## 使用互斥锁

```go
package main

import (
	"context"
	"log"

	"github.com/MaricoHan/redisson"
	"github.com/go-redis/redis/v8"
)

func main() {
	// 1.初始化 redis 连接
	client := redis.NewClient(&redis.Options{Addr: ":6379"})

	// 2.基于该连接，初始化一个 redisson
	r := redisson.New(context.Background(), client)

	// 3.初始化一把锁
	mutex := r.NewMutex("mutexKey")

	// 4.上锁
	err := mutex.Lock(context.Background())
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("lock successfully!")

	// 5.你需要处理的任务
	// ...

	// 6.解锁
	err = mutex.Unlock(context.Background())
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("unlock successfully!")
}
```

## 使用读写锁

### 写锁

```go
package main

import (
	"context"
	"log"

	"github.com/MaricoHan/redisson"
	"github.com/go-redis/redis/v8"
)

func main() {
	// 1.初始化 redis 连接
	client := redis.NewClient(&redis.Options{Addr: ":6379"})

	// 2.基于该连接，初始化一个 redisson
	r := redisson.New(context.Background(), client)

	// 3.初始化一把锁
	mutex := r.NewRWMutex("rwMutexKey")

	// 4.上写锁
	err := mutex.Lock(context.Background())
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("lock successfully!")

	// 5.你需要处理的任务
	// ...

	// 6.解锁
	err = mutex.Unlock(context.Background())
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("unlock successfully!")
}
```

### 读锁

```go
package main

import (
	"context"
	"log"

	"github.com/MaricoHan/redisson"
	"github.com/go-redis/redis/v8"
)

func main() {
	// 1.初始化 redis 连接
	client := redis.NewClient(&redis.Options{Addr: ":6379"})

	// 2.基于该连接，初始化一个 redisson
	r := redisson.New(context.Background(),client)

	// 3.初始化一把锁
	rwMutex := r.NewRWMutex("rwMutexKey")

	// 4.上读锁
	err := rwMutex.RLock(context.Background())
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("rLock successfully!")

	// 5.你需要处理的任务
	// ...

	// 6.解锁
	err = rwMutex.Unlock(context.Background())
	if err != nil {
		log.Panic(err)
		return
	}
	log.Println("unlock successfully!")
}
```

### 可选项

可以设置以下可选项，当不设置时，默认 `WaitTimeout = 30s`、`ExpireDuration = 10s`；

```go
options := []mutex.Option{
    mutex.WithWaitTimeout(10 * time.Second),      // 抢锁的最长等待时间
    mutex.WithExpireDuration(20 * time.Second),   // 持有锁的过期时间
}
rwMutex := r.NewRWMutex("rwMutexKey", options...) // 互斥锁同理
```
