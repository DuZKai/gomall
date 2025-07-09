package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gomall/app/seckill/biz/dal/redis"
	"net/http"
)

func SeckillStatusHandler(c *gin.Context) {
	userID := c.Query("user_id")
	activityID := c.Query("activity_id")
	if userID == "" || activityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
		return
	}

	tokenKey := fmt.Sprintf("seckill_token:%s:%s", userID, activityID)
	failKey := fmt.Sprintf("seckill_fail:%s:%s", userID, activityID)

	ctx := c.Request.Context()

	// 先查 token 是否存在
	token, err := redis.RedisClient.Get(ctx, tokenKey).Result()
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "token": token})
		return
	}

	// 查是否存在失败标记
	failFlag, _ := redis.RedisClient.Exists(ctx, failKey).Result()
	if failFlag > 0 {
		c.JSON(http.StatusOK, gin.H{"status": "FAIL"})
		return
	}

	// 否则视为排队中
	c.JSON(http.StatusOK, gin.H{"status": "PENDING"})
}
