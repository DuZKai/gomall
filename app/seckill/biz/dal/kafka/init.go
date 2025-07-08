package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"gomall/app/seckill/biz/service"
	"gomall/app/seckill/conf"
	"log"
)

var KafkaProducer sarama.AsyncProducer
var ConsumerGroup sarama.ConsumerGroup

func InitKafkaProducer() {
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
	KafkaProducer = producer
}

func Init() {
	InitKafkaProducer()
}

func InitKafkaConsumerGroup(ctx context.Context) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	groupID := "seckill_consumer_group"

	group, err := sarama.NewConsumerGroup(conf.GetConf().Kafka.Address, groupID, config)
	if err != nil {
		panic(err)
	}

	ConsumerGroup = group

	log.Println("[Consumer] Started. Waiting for messages...")
	go func() {
		for {
			if err := ConsumerGroup.Consume(ctx, []string{conf.GetConf().Kafka.Topic}, &service.SeckillConsumer{}); err != nil {
				log.Printf("[Consumer] Error: %v\n", err)
			}
			// 若 ctx 被 cancel，退出
			if ctx.Err() != nil {
				log.Println("[Consumer] Context cancelled, stopping.")
				return
			}
		}
	}()
}
