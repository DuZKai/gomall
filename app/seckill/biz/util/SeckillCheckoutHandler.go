package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gomall/app/seckill/biz/dal/mysql"
	rc "gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"net/http"
	"time"
)

func SeckillCheckoutHandler(c *gin.Context) {
	var req model.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Redis 校验 token 是否存在
	tokenKey := fmt.Sprintf("seckill:token:%s:%s", req.ActivityID, req.UserID)
	ctx := c.Request.Context()
	_, err := rc.RedisClient.Get(ctx, tokenKey).Result()
	if err == redis.Nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "token expired or invalid"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "redis error"})
		return
	}

	// 开启事务下单
	db := mysql.DB // 获取数据库连接
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "start transaction failed"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "panic rollback"})
		}
	}()

	// 插入订单
	err = tx.Create(&model.Order{
		UserID:     req.UserID,
		ActivityID: req.ActivityID,
		Status:     "INIT",
		CreateTime: time.Now(),
	}).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create order failed"})
		return
	}

	// Redis 删除 token
	err = rc.RedisClient.Del(ctx, tokenKey).Err()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete token failed"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
