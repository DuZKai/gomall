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
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	BaseTokenRate     = 1200
	TokenBucketFactor = 5
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
	captcha := req.Captcha

	if userID == "" || activityID == "" || captcha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}
	log.Printf("[SeckillRequest] user=%s, activity=%s, priority=%d, captcha=%s", userID, activityID, req.Priority, captcha)

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
	// 4. 判断库存是否充足
	stockKey := fmt.Sprintf("seckill:stock:%s", activityID)
	stockStr, err := rc.RedisClient.Get(ctx, stockKey).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 字符串转整数
	stockNum, err := strconv.Atoi(stockStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid stock format"})
		return
	}

	// 判断库存是否为0
	if stockNum <= 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "activity stock is empty"})
		return
	}

	// 校验验证码，屏蔽使得程序正常运行
	// if !verifyCaptcha(captcha) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid captcha"})
	// 	return
	// }
	// 判断是否同一个人，5秒只能一次
	// userKey := fmt.Sprintf("seckill:user:%s:%s", activityID, userID)
	// ok, err := rc.RedisClient.SetNX(ctx, userKey, 1, time.Second*5).Result()
	// if err != nil {
	// 	panic(err)
	// }
	// if !ok {
	// 	// 用户已经请求过了
	// 	c.JSON(http.StatusTooManyRequests, gin.H{"error": "You have already requested this activity"})
	// 	return
	// }

	// 第二步：限流判断 - 根据优先级选择限流策略
	if req.Priority == 1 { // VIP用户
		// 使用Sentinel匀速排队限流（资源名seckill_vip）
		entry, blockErr := api.Entry("seckill_vip", api.WithTrafficType(base.Inbound))
		if blockErr != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests - rate limited (sentinel)"})
			return
		}
		defer entry.Exit()
	} else if req.Priority == 0 { // 普通用户
		// 使用Redis令牌桶限流
		// rate: 令牌生成速率（建议=预期QPS×1.2）
		// capacity: 桶容量（建议=库存×5）
		if !AllowByTokenBucket(activityID, BaseTokenRate, stockNum*TokenBucketFactor) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests - rate limited (token bucket)"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid priority"})
		return
	}

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

	select {
	case kafka.KafkaProducer.Input() <- msg:
		log.Printf("[Kafka] 投递成功: user=%s, activity=%s", userID, activityID)

		msgTimeKey := fmt.Sprintf("seckill:msg_time:%s:%s", activityID, userID)
		err := rc.RedisClient.Set(ctx, msgTimeKey, time.Now().Unix(), 2*time.Minute).Err()
		if err != nil {
			log.Printf("[Redis] 设置 msg_time 失败: %v", err)
			// 如果写入失败可以继续处理（非致命），也可以根据策略决定是否回滚 Kafka 投递
		}
		c.JSON(http.StatusOK, gin.H{"message": "Request accepted, queuing..."})

	default:
		log.Printf("[Kafka] 投递失败: user=%s, activity=%s", userID, activityID)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "server busy, try again later"})
		return
	}
	// 查看容器内是否有消息
	// docker exec -it kafka bash
	// kafka-console-consumer.sh \
	//  --bootstrap-server localhost:9092 \
	//  --topic seckill_requests \
	//  --from-beginning

	log.Printf("[SeckillRequest][OK] user=%s 请求进入队列", userID)
}
