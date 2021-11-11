package redis

import (
	"log"
	"plog_gateway/db"
)

func Publish(queueKey string, message string) {
	rdb := db.DialRedis()
	defer rdb.Close()
	r, err := rdb.Publish(queueKey, message).Result()
	if err != nil {
		log.Println(err, r)
	}
}
