package asynq

import (
	"context"
	"gomall/app/seckill/biz/dal/redis"
	"time"
)

func ShutdownAll() {
	ctx := context.Background()
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	AsynqScheduler.Shutdown()
	AsynqServer.Shutdown()
	if err := AsynqClient.Close(); err != nil {
		panic(err)
	}
	// 清理所有的 asynq 相关的 Redis 键
	keys, err := redis.RedisClient.Keys(ctx, "asynq:*").Result()
	if err == nil && len(keys) > 0 {
		redis.RedisClient.Del(ctx, keys...)
	}
}
