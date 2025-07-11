package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gomall/app/seckill/biz/dal/kafka"
	rc "gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"gomall/app/seckill/conf"
	"net/http"
	"time"
)

func SeckillRequestHandler(c *gin.Context) {
	ctx := context.Background()
	// 第一步: 验证资格接口
	var req model.SeckillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	userID := req.UserID
	activityID := req.ActivityID
	captcha := req.ActivityID

	if userID == "" || activityID == "" || captcha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	// 判断活动是否开始
	activityKey := fmt.Sprintf("seckill:activity:%s", activityID)
	val, err := rc.RedisClient.Get(ctx, activityKey).Result()
	if err == redis.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "token expired or invalid"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
		return
	}

	// 1. 反序列化活动信息
	var activity model.Activity
	if err := json.Unmarshal([]byte(val), &activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "activity parse failed"})
		return
	}

	// 2. 获取时间
	now := time.Now().Unix()
	startTime := activity.StartTime
	endTime := activity.EndTime

	// 3. 判断是否在活动时间内
	if now < startTime {
		c.JSON(http.StatusForbidden, gin.H{"error": "activity has not started"})
		return
	}
	if now >= endTime {
		c.JSON(http.StatusForbidden, gin.H{"error": "activity has ended"})
		return
	}

	// if !verifyCaptcha(captcha) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid captcha"})
	// 	return
	// }
	// TODO 事实应该判断是否同一个人

	// 第二步：限流判断
	// TODO：按照消息队列长度限流
	entry, blockErr := api.Entry(conf.GetConf().Kafka.Topic, api.WithTrafficType(base.Inbound))
	if blockErr != nil {
		// 被限流
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests - rate limited"})
		return
	}
	defer entry.Exit()

	// 第三步：构造 Kafka 消息并异步发送
	msgBody := map[string]string{
		"user_id":     userID,
		"activity_id": activityID,
		"ts":          time.Now().Format(time.RFC3339),
	}
	jsonBytes, _ := json.Marshal(msgBody)

	msg := &sarama.ProducerMessage{
		Topic: conf.GetConf().Kafka.Topic,
		Key:   sarama.StringEncoder(activityID), // 按活动 ID 分区
		Value: sarama.ByteEncoder(jsonBytes),
	}

	kafka.KafkaProducer.Input() <- msg
	// 查看容器内是否有消息
	// docker exec -it kafka bash
	// kafka-console-consumer.sh \
	//  --bootstrap-server localhost:9092 \
	//  --topic seckill_requests \
	//  --from-beginning

	c.JSON(http.StatusOK, gin.H{"message": "Request accepted, queuing..."})

}
