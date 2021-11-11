package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// 同步消息模式,把消息发送后等待结果返回
func SyncProducerWithMessage(address []string, topic, value string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Timeout = time.Second * 5

	producer, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		logrus.Errorf("Failed to create sync producer: %v", err)
		return err
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		logrus.Errorf("Failed to send message: %v", err)
		return err
	}
	return nil
}

// 异步消息模式,把消息发送后不等待结果返回
func AsyncProducerWithMessage(address []string, topic, value string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	// 等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Version = sarama.V0_10_0_1

	producer, err := sarama.NewAsyncProducer(address, config)
	if err != nil {
		logrus.Errorf("Failed to create async producer: %v", err)
		return err
	}
	defer producer.AsyncClose()
	producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.ByteEncoder(value)}

	select {
	case <-producer.Successes():
		logrus.Info("go go go!")
	case fail := <-producer.Errors():
		logrus.Infof("err:%v", fail.Err)
	default:
		logrus.Info("go ?")
	}
	return nil
}
