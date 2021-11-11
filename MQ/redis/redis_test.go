package redis

import (
	"testing"
	"time"
)

func TestRedisMQ(t *testing.T) {
	go Consume("TestQueue")
	time.Sleep(time.Second * 2)
	Publish("TestQueue", "pp")

	Publish("TestQueue", "apale")
}
