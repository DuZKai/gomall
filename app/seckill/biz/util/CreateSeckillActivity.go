package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gomall/app/seckill/biz/dal/mysql"
	rc "gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"net/http"
	"time"
)

// 创建活动接口
func CreateSeckillActivity(c *gin.Context) {
	var req model.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 1. 模拟入库（此处应调用实际 db.InsertActivity）
	activity := model.Activity{
		ID:        req.ActivityID,
		ProductID: req.ProductID,
		Stock:     req.Stock,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Remark:    req.Remark,
		CreateAt:  time.Now().Unix(),
	}

	if err := mysql.DB.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "活动入库失败"})
		return
	}

	// 2. 写入 Redis 缓存（活动信息 + 初始库存）
	err := cacheSeckillActivity(activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缓存写入失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "创建成功",
		"activity_id": activity.ID,
	})
}

// 写缓存逻辑
func cacheSeckillActivity(a model.Activity) error {
	ctx := context.Background()
	now := time.Now().Unix()
	expireSeconds := (a.EndTime-a.StartTime)*2 + (a.StartTime - now)
	if expireSeconds <= 0 {
		expireSeconds = 600 // 给个最小兜底值，例如10分钟
	}

	expire := time.Duration(expireSeconds) * time.Second

	// 活动信息缓存
	activityKey := fmt.Sprintf("seckill:activity:%d", a.ID)
	activityJson, _ := json.Marshal(a)
	err := rc.RedisClient.Set(ctx, activityKey, activityJson, expire).Err()
	if err != nil {
		return err
	}

	// 库存初始化缓存
	stockKey := fmt.Sprintf("seckill_stock:%d", a.ID)
	err = rc.RedisClient.Set(ctx, stockKey, a.Stock, expire).Err()
	return err
}
