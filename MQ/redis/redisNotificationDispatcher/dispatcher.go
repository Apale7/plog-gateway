package redisNotificationDispatcher

import (
	"fmt"
	"plog_gateway/MQ/redis/redis_service"
	"plog_gateway/utils"
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
)

type Subscription struct {
	IsPattern bool
	Channel   string
}

type Dispatcher struct {
	conn        *redis.PubSub
	channels    map[Subscription]*_Subscribers
	mtxChannels sync.RWMutex
	exit        bool
}

type Subscriber interface {
	// 目前的实现偏向性能，但高频推送的时候可能会有少量重复的消息，对于通知推送来说，重复消息一般是可以接受的
	OnNotifyMessage(msg *redis.Message, skipped int64)
	OnNotifyError(error)
}

type _Subscribers struct {
	v map[Subscriber]*_Subscriber
	m sync.RWMutex
}

type _Subscriber struct {
	subscriber Subscriber
	event      *utils.Event
	lastMsg    *redis.Message
	skipped    int64
	chErr      chan error
	exitEvent  *utils.Event
}

var (
	subscribeDispatchers = make(map[string]*Dispatcher)
	defaultQueueSize     = 50
)

// 设置默认队列长度
func SetDefaultQueueSize(size int) {
	defaultQueueSize = size
}

// 获取一个Dispatcher
func Get(configName string) *Dispatcher {
	if configName == `` {
		configName = `default`
	}
	dispatcher, ok := subscribeDispatchers[configName]
	if !ok {
		dispatcher = &Dispatcher{
			conn:     redis_service.GetGoRedis(0, configName).Subscribe(),
			channels: make(map[Subscription]*_Subscribers),
		}
		subscribeDispatchers[configName] = dispatcher
		go dispatcher.run()
	}
	return dispatcher
}

// 订阅指定主题
func (me *Dispatcher) Subscribe(subscriber Subscriber, channels ...string) error {
	return me.doSubscribe(subscriber, false, channels...)
}

// 订阅指定格式的主题
func (me *Dispatcher) PSubscribe(subscriber Subscriber, channels ...string) error {
	return me.doSubscribe(subscriber, true, channels...)
}

func (me *Dispatcher) doSubscribe(subscriber Subscriber, isPattern bool, channels ...string) error {
	subscription := Subscription{
		IsPattern: isPattern,
	}

	for _, c := range channels {
		subscription.Channel = c
		me.mtxChannels.Lock()

		subscribers, ok := me.channels[subscription]
		if !ok { // 新的订阅
			subscribers = &_Subscribers{
				v: make(map[Subscriber]*_Subscriber),
			}
		}

		subscribers.m.Lock()
		wrappedSubscriber, ok2 := subscribers.v[subscriber]
		if !ok2 { // 这个subscriber之前没有订阅过该主题,因此为他创建处理队列
			wrappedSubscriber = &_Subscriber{
				subscriber: subscriber,
				event:      utils.NewEvent(),
				chErr:      make(chan error),
				exitEvent:  utils.NewEvent(),
			}

			subscribers.v[subscriber] = wrappedSubscriber

			go wrappedSubscriber.Run()
		}
		subscribers.m.Unlock()

		if !ok { // 最后再放到channel里面去,避免一订阅就收到消息,如果这时候订阅通道列表已经有了,但是订阅者列表是空的,会取消订阅通道
			me.channels[subscription] = subscribers
			var err error
			if isPattern {
				err = me.conn.PSubscribe(channels...)
			} else {
				err = me.conn.Subscribe(channels...)
			}
			if err != nil {
				return fmt.Errorf("Failed to subscribe %+v", err)
			}
		}

		me.mtxChannels.Unlock()
	}
	return nil
}

func (me *Dispatcher) run() {
	type LastMessage struct {
		msg         *redis.Message
		subscribers *_Subscribers
	}

	for !me.exit {

		v, err := me.conn.Receive() // 上一次把消息读完了,这次等待新消息
		if me.exit {
			return
		}

		if err != nil {
			logrus.Errorf(`[GoRedisSubscribeDispatcher.Run] Failed to receive: %s`, err.Error())

			me.mtxChannels.RLock()
			for _, c := range me.channels {
				c.m.RLock()
				for _, subscriber := range c.v {
					if len(subscriber.chErr) <= cap(subscriber.chErr) {
						subscriber.chErr <- err
					}
				}
				c.m.RUnlock()
			}
			me.mtxChannels.RUnlock()
			continue
		}

		switch msg := v.(type) {
		case *redis.Message:
			pattern := msg.Pattern

			if pattern == `` {
				me.dispatch(Subscription{
					IsPattern: false,
					Channel:   msg.Channel,
				}, msg)
			} else {
				me.dispatch(Subscription{
					IsPattern: true,
					Channel:   pattern,
				}, msg)
			}
		case *redis.Subscription:
		}
	}
}

func (me *Dispatcher) dispatch(subscription Subscription, msg *redis.Message) {
	me.mtxChannels.RLock()
	subscribers, ok := me.channels[subscription]
	me.mtxChannels.RUnlock()
	if !ok {
		logrus.Errorf(`[redisNotificationManager.Dispatch] subscription %+v not found!`, subscription)
		return
	}

	subscribers.m.RLock()
	for _, subscriber := range subscribers.v {
		subscriber.lastMsg = msg
		atomic.AddInt64(&subscriber.skipped, 1)
		subscriber.event.Set()
	}
	subscribers.m.RUnlock()
}

// 暂时单例应该不会close
func (me *Dispatcher) close() {
	me.conn.Close()
}

func (me *_Subscriber) Run() {
	for {
		select {
		case <-me.event.Done():
			me.event.Unset()
			msg := me.lastMsg
			skipped := atomic.SwapInt64(&me.skipped, 0)
			me.subscriber.OnNotifyMessage(msg, skipped-1)
		case err := <-me.chErr:
			me.subscriber.OnNotifyError(err)
		case <-me.exitEvent.Done():
			logrus.Debugf(`Subscriber %+v don`, me)
			return
		}
	}
}

func (me *_Subscriber) close() {
	me.exitEvent.Set()
}
