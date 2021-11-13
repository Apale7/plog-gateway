package redisSubscriptionDispatcher

import "time"

// 公有定义
type (
	NewConnParams struct {
		RedisConf     string
		BatchSize     int           // 消息是批量接收批量派发的,这里指定一批多少条消息,默认1000
		BatchInterval time.Duration // 批量接收消息队列的超时时间,默认1ms
	}
)

// 私有定义
type ()
