package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gomall/app/seckill/biz/dal/mysql"
	rc "gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func SeckillCheckoutHandler(c *gin.Context) {
	log.Println("[CheckoutHandler] Processing checkout request")
	var req model.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Redis 校验 token 是否存在
	tokenKey := fmt.Sprintf("seckill:token:valid:%s:%s", req.ActivityID, req.UserID)
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

	// 判断是否重复下单（幂等处理）
	var existing model.Order
	if err := tx.Where("user_id = ? AND activity_id = ?", req.UserID, req.ActivityID).First(&existing).Error; err == nil {
		tx.Rollback()
		log.Println("[CheckoutHandler] Duplicate order attempt detected")
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate order"})
		return
	}

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

	res := tx.Model(&model.ActivityStock{}).
		Where("activity_id = ? AND stock > 0", req.ActivityID).
		Update("stock", gorm.Expr("stock - ?", 1))
	if res.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "deduct stock failed"})
		return
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "out of stock"})
		return
	}

	// 提交事务并检查错误
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "commit transaction failed"})
		return
	}

	// 原子性校验 + 删除 token（Lua）
	luaScript := `
		if redis.call("GET", KEYS[1]) then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`
	delRes, err := rc.RedisClient.Eval(ctx, luaScript, []string{tokenKey}).Result()
	if err != nil || delRes.(int64) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "token expired or already used"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
	log.Println("[SeckillCheckoutHandler] Checkout success")
}
