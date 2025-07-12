package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	rc "gomall/app/seckill/biz/dal/redis"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SeckillStatusHandler(c *gin.Context) {
	userID := c.Query("user_id")
	activityID := c.Query("activity_id")
	if userID == "" || activityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
		return
	}

	tokenKey := fmt.Sprintf("seckill:token:valid:%s:%s", activityID, userID)
	failKey := fmt.Sprintf("seckill:fail:%s:%s", activityID, userID)

	ctx := c.Request.Context()

	// 先查 token 是否存在
	token, err := rc.RedisClient.Get(ctx, tokenKey).Result()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "token": token})
		return
	}
	if err != redis.Nil {
		log.Printf("[StatusHandler][RedisError] get tokenKey failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal redis error"})
		return
	}

	// 查是否存在失败标记
	failFlag, _ := rc.RedisClient.Exists(ctx, failKey).Result()
	if failFlag > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "FAIL"})
		return
	}

	// 如果 Redis 无 token 且无 failKey，则查询 msg_time 超时时间
	msgTimeKey := fmt.Sprintf("seckill:msg_time:%s:%s", activityID, userID)
	msgTimeStr, err := rc.RedisClient.Get(ctx, msgTimeKey).Result()
	if err == redis.Nil {
		c.JSON(http.StatusOK, gin.H{"status": "NOT_FOUND"})
		return
	} else if err != nil {
		log.Printf("[StatusHandler] Redis error on msg_time: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
		return
	}

	msgTime, _ := strconv.ParseInt(msgTimeStr, 10, 64)
	if time.Now().Unix()-msgTime > 60 {
		c.JSON(http.StatusOK, gin.H{"status": "FAIL", "reason": "timeout"})
		return
	}

	// 否则视为排队中
	c.JSON(http.StatusOK, gin.H{"status": "PENDING"})
}
