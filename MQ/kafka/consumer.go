package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

func ConsumerWithMessage(address []string, topic string, partition int32, offset int64) error {
	fmt.Println("Consumer")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(address, config)
	if err != nil {
		logrus.Errorf("Failed to create consumer: %v", err)
		return err
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
	if err != nil {
		logrus.Errorf("Failed to careate consume partition: %v", err)
		return err
	}
	defer partitionConsumer.Close()
	for msg := range partitionConsumer.Messages() {
		fmt.Printf("partition:%d offset:%d key:%s val:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
	return nil
}
