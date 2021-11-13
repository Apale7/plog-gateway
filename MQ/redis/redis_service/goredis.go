package redis_service

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

// RWMutexRedisPools 替代原先单纯的map，可以加锁避免同时读写错误
type RWMutexRedisPools struct {
	goRedisPools map[string]*redis.Client
	sync.RWMutex
}

var goRWMutexRedisPools = RWMutexRedisPools{goRedisPools: make(map[string]*redis.Client)}

// GetInstance 获取redis连接池
func GetGoRedis(db int, section ...string) *redis.Client {
	var sectionName = "default"
	if len(section) > 0 {
		sectionName = section[0]
	}
	goRWMutexRedisPools.RLock()
	if result, ok := goRWMutexRedisPools.goRedisPools[sectionName]; ok {
		goRWMutexRedisPools.RUnlock()
		return result
	}
	goRWMutexRedisPools.RUnlock()

	conf, ok := confs[sectionName]

	if !ok {
		logrus.Errorf(`go-redis config "%s" not registered`, sectionName)
		return nil
	}

	result := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(`%s:%d`, conf.host, conf.port),
		Password: conf.password,
		DB:       db,
		PoolSize: conf.poolMaxConns,
	})
	goRWMutexRedisPools.Lock()
	goRWMutexRedisPools.goRedisPools[sectionName] = result
	goRWMutexRedisPools.Unlock()
	return result
}

// SubscribeGoRedis 获取一个发布/订阅器
func SubscribeGoRedis(channels ...string) *redis.PubSub {
	return GetGoRedis(0).Subscribe(channels...)
}

// PSubscribeGoRedis 获取一个发布/订阅器
func PSubscribeGoRedis(patterns ...string) *redis.PubSub {
	return GetGoRedis(0).PSubscribe(patterns...)
}
