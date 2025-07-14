package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"gomall/app/seckill/biz/service"
	"gomall/app/seckill/conf"
	"log"
	"time"
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

func CreateTopic(brokers []string, topic string, numPartitions int32, replicationFactor int16) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return err
	}
	defer func(admin sarama.ClusterAdmin) {
		err := admin.Close()
		if err != nil {
			log.Printf("Error closing cluster admin: %v", err)
			panic(err)
		}
	}(admin)

	// 创建 topic
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		// 如果 topic 已经存在，可以忽略
		if err.(*sarama.TopicError).Err == sarama.ErrTopicAlreadyExists {
			log.Printf("Topic %s already exists", topic)
		} else {
			return err
		}
	} else {
		log.Printf("Topic %s created successfully", topic)
	}
	return nil
}

func InitKafkaConsumerGroup(brokers []string, groupID string, topic string) error {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}
	ConsumerGroup = group

	handler := &service.SeckillConsumer{}

	ctx := context.Background()
	go func() {
		log.Printf("[Consumer] Started.")
		for {
			if err := group.Consume(ctx, []string{topic}, handler); err != nil {
				log.Printf("[Consumer] Error: %v", err)
				time.Sleep(time.Second)
			}
			if ctx.Err() != nil {
				log.Println("[Consumer] Context cancelled, exiting.")
				return
			}
		}
	}()
	return nil
}

func InitKafkaConsumer() {
	topic := conf.GetConf().Kafka.Topic
	numPartitions := int32(16)
	replicationFactor := int16(1)
	groupID := "seckill_consumer_group"

	// 创建 topic
	if err := CreateTopic(conf.GetConf().Kafka.Address, topic, numPartitions, replicationFactor); err != nil {
		log.Fatalf("CreateTopic failed: %v", err)
	}

	// 初始化消费者组，启动多个消费者实例
	err := InitKafkaConsumerGroup(conf.GetConf().Kafka.Address, groupID, topic)
	if err != nil {
		log.Fatalf("InitKafkaConsumerGroup failed: %v", err)
	}
}
