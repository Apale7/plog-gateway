package redis

import (
	"log"
	"plog_gateway/db"
)

func Consume(queueKey string) {
	rdb := db.DialRedis()
	pubSub := rdb.Subscribe(queueKey)
	ch := pubSub.Channel()
	for msg := range ch {
		log.Printf("消费到数据，channel:%s；message:%s\n", msg.Channel, msg.Payload)
	}
}
