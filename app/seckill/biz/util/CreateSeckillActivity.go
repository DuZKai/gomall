package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gomall/app/seckill/biz/dal"
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

	// 生成一个雪花ID
	id := dal.Node.Generate()
	const timeLayout = "2006-01-02 15:04:05"
	startTime, err := time.ParseInLocation(timeLayout, req.StartTime, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间格式错误"})
		return
	}

	endTime, err := time.ParseInLocation(timeLayout, req.EndTime, time.Local)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间格式错误"})
		return
	}

	// 1. 入库
	activity := model.Activity{
		ID:         id.String(),
		ActivityID: req.ActivityID,
		ProductID:  req.ProductID,
		Stock:      req.Stock,
		StartTime:  startTime.Unix(),
		EndTime:    endTime.Unix(),
		Remark:     req.Remark,
		CreateAt:   time.Now().Unix(),
	}

	if err := mysql.DB.Create(&activity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "活动入库失败"})
		return
	}

	// 2. 写入 Redis 缓存（活动信息 + 初始库存）
	err = cacheSeckillActivity(activity, startTime.Unix(), endTime.Unix())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缓存写入失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"id":      activity.ID,
	})
}

// 写缓存逻辑
func cacheSeckillActivity(a model.Activity, startTime int64, endTime int64) error {
	ctx := context.Background()
	now := time.Now().Unix()
	expireSeconds := (endTime-startTime)*2 + (startTime - now)
	if expireSeconds <= 0 {
		expireSeconds = 600 // 给个最小兜底值，例如10分钟
	}

	expire := time.Duration(expireSeconds) * time.Second

	// 活动信息缓存
	activityKey := fmt.Sprintf("seckill:activity:%s", a.ActivityID)
	activityJson, _ := json.Marshal(a)
	err := rc.RedisClient.Set(ctx, activityKey, activityJson, expire).Err()
	if err != nil {
		return err
	}

	// 库存初始化缓存
	stockKey := fmt.Sprintf("seckill:stock:%s", a.ActivityID)
	err = rc.RedisClient.Set(ctx, stockKey, a.Stock, expire).Err()

	// 令牌桶初始化
	key := fmt.Sprintf("seckill:token:bucket:%s", a.ActivityID)
	now = time.Now().UnixNano() / 1e6 // 毫秒
	err = rc.RedisClient.HSet(ctx, key, map[string]interface{}{
		"last_mill_second": now,
		"tokens":           2000,
	}).Err()
	if err != nil {
		return err
	}

	// 设置 key 过期时间
	err = rc.RedisClient.Expire(ctx, key, expire).Err()
	if err != nil {
		return err
	}
	return err
}
