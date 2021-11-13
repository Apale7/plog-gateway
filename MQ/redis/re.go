package redis

import (
	"fmt"
	"sync"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

var (
	confs map[string]_rdConf
)

type _rdConf struct {
	db           int
	host         int
	port         int
	password     string
	poolMaxIdle  int
	poolMaxConns int
}

type RWMutexRedisPools struct {
	goRedisPools map[string]*redis.Client
	sync.RWMutex
}

var goRWMutexRedisPools = RWMutexRedisPools{
	goRedisPools: make(map[string]*redis.Client),
}

func init() {
	confs = make(map[string]_rdConf)
}

func getConf(sectionName string, db int) (_rdConf, bool) {
	cf, ok := confs[sectionName]
	if ok {
		cf.db = db
	}
	return cf, ok
}

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

	conf, ok := getConf(sectionName, db)

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
	_, err := result.Ping().Result()
	if err != nil {
		logrus.Errorf(`Create redis "%s" failed`, sectionName)
		return nil
	}
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
