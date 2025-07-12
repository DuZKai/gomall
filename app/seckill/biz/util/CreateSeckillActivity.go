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
	"gomall/app/seckill/config"
	"log"
	"net/http"
	"time"
)

// 创建活动接口
func CreateSeckillActivity(c *gin.Context) {
	var req model.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[CreateActivity] 参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 生成一个雪花ID
	id := dal.Node.Generate()
	const timeLayout = "2006-01-02 15:04:05"
	startTime, err := time.ParseInLocation(timeLayout, req.StartTime, time.Local)
	if err != nil {
		log.Printf("[CreateActivity] 开始时间解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "开始时间格式错误"})
		return
	}

	endTime, err := time.ParseInLocation(timeLayout, req.EndTime, time.Local)
	if err != nil {
		log.Printf("[CreateActivity] 结束时间解析失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "结束时间格式错误"})
		return
	}

	db := mysql.DB // 获取数据库连接
	var exist model.Activity
	if err := db.Where("activity_id = ?", req.ActivityID).First(&exist).Error; err == nil {
		log.Printf("[CreateActivity] 活动已存在: activity_id=%s", req.ActivityID)
		c.JSON(http.StatusConflict, gin.H{"error": "活动已存在"})
		return
	}

	// 1. 入库
	tx := db.Begin()
	if tx.Error != nil {
		log.Printf("[CreateActivity] 开启事务失败: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "开启事务失败"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("[CreateActivity] 回滚失败: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "回滚事务失败"})
		}
	}()

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

	if err := tx.Create(&activity).Error; err != nil {
		tx.Rollback()
		log.Printf("[CreateActivity] 活动入库失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "活动入库失败"})
		return
	}

	// 2. 插入库存表
	stock := model.ActivityStock{
		ActivityID: req.ActivityID,
		Stock:      req.Stock,
	}
	if err := tx.Create(&stock).Error; err != nil {
		tx.Rollback()
		log.Printf("[CreateActivity] 库存入库失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "库存入库失败"})
		return
	}

	// 3. 提交事务
	if err := tx.Commit().Error; err != nil {
		log.Printf("[CreateActivity] 事务提交失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "事务提交失败"})
		return
	}

	// 2. 写入 Redis 缓存（活动信息 + 初始库存）
	err = cacheSeckillActivity(activity, startTime.Unix(), endTime.Unix())
	if err != nil {
		log.Printf("[CreateActivity] Redis 缓存失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "缓存写入失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"id":      activity.ID,
	})
	log.Printf("[CreateActivity] 活动创建成功: activity_id=%s stock=%d", activity.ActivityID, activity.Stock)
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
		"tokens":           a.Stock * int64(config.AppConfig.CapacityFactor), // 初始令牌数为库存的5倍
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
