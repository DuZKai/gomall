package util

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"gomall/app/seckill/biz/dal/kafka"
	"net/http"
	"time"
)

type SeckillRequest struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	Captcha    string `json:"captcha"`
}

func SeckillRequestHandler(c *gin.Context) {
	// 第一步: 验证资格接口
	var req SeckillRequest
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

	// if !verifyCaptcha(captcha) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid captcha"})
	// 	return
	// }
	// TODO 事实应该判断是否同一个人

	// 第二步：限流判断
	// TODO：按照消息队列长度限流
	entry, blockErr := api.Entry("seckill_request", api.WithTrafficType(base.Inbound))
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
		Topic: "seckill_requests",
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
