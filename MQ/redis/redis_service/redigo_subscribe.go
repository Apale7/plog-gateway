package redis_service

import (
	"fmt"
	"net"
	"plog_gateway/utils"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/gomodule/redigo/redis"
)

type SubConn struct {
	connectionName   string
	channels         []interface{}
	containsPatterns bool
	lastConn         redis.Conn
	lastErr          error
}

func (me *SubConn) get() (redis.Conn, error) {
	if me.lastConn == nil || me.lastConn.Err() != nil {
		cf, _ := getConf(me.connectionName, 0)
		me.lastConn, me.lastErr = redis.Dial(`tcp`, fmt.Sprintf(`%s:%d`, cf.host, cf.port), redis.DialPassword(cf.password),
			redis.DialConnectTimeout(time.Second*10),
			redis.DialWriteTimeout(time.Second*3))
		if me.lastErr != nil {
			return nil, me.lastErr
		}
		if me.containsPatterns {
			_, me.lastErr = me.lastConn.Do(`PSUBSCRIBE`, me.channels...)
		} else {
			_, me.lastErr = me.lastConn.Do(`SUBSCRIBE`, me.channels...)
		}
	}
	return me.lastConn, me.lastErr
}

func (me *SubConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn, err := me.get()
	if err != nil {
		return nil, err
	}
	return conn.Do(commandName, args...)
}

func (me *SubConn) Close() error {
	if me.lastConn != nil {
		return me.lastConn.Close()
	}
	return nil
}

func (me *SubConn) Err() error {
	if me.lastConn != nil {
		return me.lastConn.Err()
	}
	return me.lastErr
}

func (me *SubConn) Flush() error {
	if me.lastConn != nil {
		return me.lastConn.Flush()
	}
	return me.lastErr
}

func (me *SubConn) Receive() (reply interface{}, err error) {
	conn, err := me.get()
	if conn == nil {
		return nil, err
	}
	return conn.Receive()
}

func (me *SubConn) Send(commandName string, args ...interface{}) error {
	conn, err := me.get()
	if conn == nil {
		return err
	}
	return conn.Send(commandName, args...)
}

func (me *SubConn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	conn, err := me.get()
	if conn == nil {
		return nil, err
	}
	return conn.(redis.ConnWithTimeout).DoWithTimeout(timeout, commandName, args...)
}

func (me *SubConn) ReceiveWithTimeout(timeout time.Duration) (reply interface{}, err error) {
	conn, err := me.get()
	if conn == nil {
		return nil, err
	}
	return conn.(redis.ConnWithTimeout).ReceiveWithTimeout(timeout)
}

// PSubscribe 自动重连、自动重新订阅
func PSubscribe(channels ...string) (*redis.PubSubConn, error) {
	var ichannels []interface{}
	linq.From(channels).ToSlice(&ichannels)
	return &redis.PubSubConn{Conn: &SubConn{
		connectionName:   `default`,
		channels:         ichannels,
		containsPatterns: utils.FindIf(len(channels), func(i int) bool { return strings.ContainsAny(channels[i], `*?[]`) }) >= 0,
	}}, nil
}

// 检测是否有新消息,没有的话返回nil,false, 有的话返回msg,true
func Peek(conn *redis.PubSubConn) (interface{}, bool) {
	msg := conn.ReceiveWithTimeout(1) // 设置为0则永不超时
	switch msg.(type) {
	case error:
		if err, ok := msg.(net.Error); ok && err.Timeout() {
			return nil, false
		}
		return msg, true
	default:
		return msg, true
	}
}
