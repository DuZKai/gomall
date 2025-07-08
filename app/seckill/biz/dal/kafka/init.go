package kafka

import (
	"github.com/IBM/sarama"
	"gomall/app/seckill/conf"
)

var KafkaProducer sarama.AsyncProducer

func InitKafkaProducer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err := sarama.NewAsyncProducer(conf.GetConf().Kafka.Address, config)
	if err != nil {
		panic(err)
	}

	go func() {
		for err := range producer.Errors() {
			panic(err)
		}
	}()

	return producer
}
