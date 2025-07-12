package util

import (
	"context"
	"fmt"
	rc "gomall/app/seckill/biz/dal/redis"
	"time"
)

func AllowByTokenBucket(activityID string, baseRate int, baseCapacity int) bool {
	script := `
        -- KEYS[1]: token bucket key
        -- KEYS[2]: stock key
        -- ARGV[1]: 当前时间戳（毫秒）
        -- ARGV[2]: 基础速率（每秒生成多少 token）
        -- ARGV[3]: 基础桶容量
        -- 返回 1 表示允许，0 表示被限流
        
        local bucket_key = KEYS[1]
        local stock_key = KEYS[2]
        local now = tonumber(ARGV[1])
        local base_rate = tonumber(ARGV[2])
        local base_capacity = tonumber(ARGV[3])
        
        -- 获取当前库存
        local stock = tonumber(redis.call("GET", stock_key)) or 0
        
        -- 动态调整参数：基于库存计算
        local dynamic_capacity = math.min(base_capacity, stock * 5)
        local dynamic_rate = math.min(base_rate, stock * 2)
        
        -- 获取令牌桶旧值
        local bucket = redis.call("HMGET", bucket_key, "last_mill_second", "tokens")
        local last_time = tonumber(bucket[1]) or 0
        local tokens = tonumber(bucket[2]) or dynamic_capacity
        
        -- 计算时间间隔新增的 tokens
        local delta = math.max(0, now - last_time)
        local added_tokens = math.floor(delta / 1000 * dynamic_rate)
        tokens = math.min(tokens + added_tokens, dynamic_capacity)
        
        -- 判断是否有可用 token
        if tokens >= 1 then
            tokens = tokens - 1
            redis.call("HMSET", bucket_key, "tokens", tokens, "last_mill_second", now)
            redis.call("EXPIRE", bucket_key, 10)  -- 设置过期时间
            return 1
        else
            return 0
        end
    `

	bucketKey := fmt.Sprintf("seckill:token:bucket:%s", activityID)
	stockKey := fmt.Sprintf("seckill:stock:%s", activityID)
	now := time.Now().UnixNano() / 1e6 // 毫秒

	res, err := rc.RedisClient.Eval(context.Background(), script, []string{bucketKey, stockKey},
		now, baseRate, baseCapacity).Result()

	if err != nil {
		// 生产环境应记录日志而不是panic
		panic(err)
		return false
	}

	return res.(int64) == 1
}
