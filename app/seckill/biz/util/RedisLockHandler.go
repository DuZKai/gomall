package util

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gomall/app/seckill/biz/dal/mysql"
	rc "gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"net/http"
	"time"
)

func RedisLockHandler(c *gin.Context) {
	// 单位秒
	lockTTL := 30 * time.Second
	ctx := context.Background()
	activityID := c.Query("activity_id")
	if activityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing activity_id"})
		return
	}

	cacheKey := fmt.Sprintf("seckill:activity:%s", activityID)

	// 1. 尝试从 Redis 缓存读取
	val, err := rc.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，反序列化返回
		var activity model.Activity
		if err := json.Unmarshal([]byte(val), &activity); err == nil {
			fmt.Println("cache hit")
			c.JSON(http.StatusOK, gin.H{"activity": activity})
			return
		}
		// 反序列化失败，继续去数据库加载（可能缓存损坏）
	} else if err != redis.Nil {
		// Redis 读错
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// err == redis.Nil 表示缓存未命中，继续走下面逻辑

	// 2. 尝试获取分布式锁，避免缓存击穿
	lockKey := fmt.Sprintf("lock:seckill:activity:%s", activityID)
	lock, err := rc.Locker.Obtain(ctx, lockKey, lockTTL, nil)
	if err == redislock.ErrNotObtained {
		// 获取锁失败，说明已有其它线程正在加载数据库
		fmt.Println("lock held by other, wait and retry")

		val2, err2 := rc.RedisClient.Get(ctx, cacheKey).Result()
		if err2 == nil {
			var activity model.Activity
			if err := json.Unmarshal([]byte(val2), &activity); err == nil {
				fmt.Println("cache hit after wait")
				c.JSON(http.StatusOK, gin.H{"activity": activity})
				return // 成功获取到缓存数据，直接返回
			}
		}
		// 依然未命中，返回错误或者空，防止请求阻塞
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "cache miss and lock held by other, please retry later"})
		return
	} else if err != nil {
		// 获取锁时发生其他错误
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保释放锁
	defer func() {
		if err := lock.Release(ctx); err != nil {
			fmt.Printf("failed to release lock: %v\n", err)
		}
	}()

	// 开启“看门狗”续期协程
	stopWatchdog := make(chan struct{})
	go func() {
		ticker := time.NewTicker(lockTTL / 2) // 每一半 TTL 续期一次
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := lock.Refresh(ctx, lockTTL, nil)
				if err != nil {
					fmt.Printf("watchdog refresh failed: %v\n", err)
					return
				}
				fmt.Println("lock refreshed")
			case <-stopWatchdog:
				return
			}
		}
	}()
	defer close(stopWatchdog)

	time.Sleep(100 * time.Second)
	// 3. 从数据库加载数据
	var activity model.Activity
	if err := mysql.DB.WithContext(ctx).First(&activity, activityID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 写入缓存，设置过期时间，防止缓存穿透存空对象
	bytes, err := json.Marshal(activity)
	if err != nil {
		fmt.Printf("json marshal failed: %v\n", err)
		c.JSON(http.StatusOK, gin.H{"activity": activity})
		return // 缓存失败不影响正常返回
	}
	err = rc.RedisClient.Set(ctx, cacheKey, bytes, 10*time.Minute).Err()
	if err != nil {
		fmt.Printf("cache set failed: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"activity": activity})
}
