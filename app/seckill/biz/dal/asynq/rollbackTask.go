package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	rc "gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"strings"
	"time"
)

const (
	TaskRollbackStock     = "seckill:rollback"
	TaskRollbackScheduler = "seckill:rollback:scheduler"
)

func NewRollbackSchedulerTask() *asynq.Task {
	return asynq.NewTask(TaskRollbackScheduler, nil)
}

func HandleRollbackSchedulerTask(ctx context.Context, t *asynq.Task) error {
	// 获取所有 seckill:token 前缀的 key
	keys, err := rc.RedisClient.Keys(ctx, "seckill:token:valid:*:*").Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		// key = seckill:token:valid:{activityID}:{userID}
		parts := strings.Split(key, ":")
		if len(parts) != 5 {
			continue
		}
		activityID := parts[3]
		userID := parts[4]

		// 判断 key 是否存在（已过期则跳过）
		val, err := rc.RedisClient.Get(ctx, key).Result()
		if err == redis.Nil {
			// token 不存在，说明已支付
			fmt.Println("token 已被删除，无需回滚")
			return nil
		} else if err != nil {
			return err
		}

		// 2. 解析 JSON
		var token model.TokenInfo
		if err := json.Unmarshal([]byte(val), &token); err != nil {
			return err
		}

		// 3. 判断是否过期
		now := time.Now().UnixNano() // 当前时间（纳秒）
		expireAt := token.CreateTime + token.ExpireSecond*1e9

		if now < expireAt {
			fmt.Println("未过期，不回滚")
			return nil
		}

		// 回滚库存
		expireSeconds := token.ExpireSecond
		stockKey := fmt.Sprintf("seckill:stock:%s", activityID)
		failKey := fmt.Sprintf("seckill:fail:%s:%s", activityID, userID)
		luaScript := `
			if redis.call("exists", KEYS[1]) == 1 then
				redis.call("incr", KEYS[2])
				redis.call("del", KEYS[1])
				-- 创建失败标记，过期时间由 ARGV[1] 指定
				redis.call("setex", KEYS[3], tonumber(ARGV[1]), 1)
				return 1
			end
			return 0
		`
		_, err = rc.RedisClient.Eval(
			ctx,
			luaScript,
			[]string{key, stockKey, failKey},
			expireSeconds, // 传给 ARGV[1]
		).Result()
		if err != nil {
			continue
		}
		fmt.Printf("rollback success: %s - %s\n", activityID, userID)
	}
	return nil
}

func AsyncInit() {
	// 注册任务处理器
	mux := asynq.NewServeMux()
	// mux.HandleFunc(TaskRollbackStock, HandleRollbackStockTask)
	mux.HandleFunc(TaskRollbackScheduler, HandleRollbackSchedulerTask)

	// 启动异步任务处理器
	go func() {
		if err := AsynqServer.Run(mux); err != nil {
			panic(err)
		}
	}()

}
