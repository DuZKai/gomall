package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	rc "gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"log"
	"strings"
	"sync"
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
	// 使用 Redis 分布式锁防止多个任务并发扫描（锁住整个任务调度）
	lockKey := "seckill:rollback:lock"
	ok, err := rc.RedisClient.SetNX(ctx, lockKey, "1", 2*time.Minute).Result()
	if err != nil || !ok {
		fmt.Println("[Rollback] 已有任务执行中，跳过本次调度")
		return nil
	}
	defer rc.RedisClient.Del(ctx, lockKey)

	// 创建并发控制的信号量通道，限制最多 20 个 goroutine 并发处理
	semaphore := make(chan struct{}, 20)
	var wg sync.WaitGroup

	var cursor uint64
	for {
		// 分批扫描 token
		keys, nextCursor, err := rc.RedisClient.Scan(ctx, cursor, "seckill:token:valid:*:*", 100).Result()
		if err != nil {
			return fmt.Errorf("redis scan error: %v", err)
		}
		cursor = nextCursor

		for _, key := range keys {
			// 启动 goroutine 处理回滚任务
			semaphore <- struct{}{}
			wg.Add(1)
			go func(k string) {
				defer func() {
					<-semaphore
					wg.Done()
					if r := recover(); r != nil {
						log.Printf("[Panic Recover] rollback key=%s panic: %v", k, r)
					}
				}()

				parts := strings.Split(k, ":")
				if len(parts) != 5 {
					log.Printf("[Skip] 格式错误 key=%s", k)
					return
				}
				activityID := parts[3]
				userID := parts[4]

				val, err := rc.RedisClient.Get(ctx, k).Result()
				if err == redis.Nil {
					// 已支付，token 自动删除，无需回滚
					return
				} else if err != nil {
					log.Printf("[Get Error] key=%s: %v", k, err)
					return
				}

				var token model.TokenInfo
				if err := json.Unmarshal([]byte(val), &token); err != nil {
					log.Printf("[Unmarshal Error] key=%s: %v", k, err)
					return
				}

				now := time.Now().UnixNano()
				expireAt := token.CreateTime + token.ExpireSecond*1e9
				if now < expireAt {
					return
				}

				// 回滚库存 Lua 脚本
				stockKey := fmt.Sprintf("seckill:stock:%s", activityID)
				failKey := fmt.Sprintf("seckill:fail:%s:%s", activityID, userID)
				luaScript := `
					if redis.call("exists", KEYS[1]) == 1 then
						redis.call("incr", KEYS[2])
						redis.call("del", KEYS[1])
						redis.call("setex", KEYS[3], tonumber(ARGV[1]), 1)
						return 1
					end
					return 0
				`

				_, err = rc.RedisClient.Eval(
					ctx,
					luaScript,
					[]string{k, stockKey, failKey},
					token.ExpireSecond,
				).Result()
				if err != nil {
					log.Printf("[Rollback Error] activity=%s user=%s: %v", activityID, userID, err)
					return
				}
				fmt.Printf("[Rollback Success] activity=%s user=%s\n", activityID, userID)
			}(key)
		}

		if cursor == 0 {
			break
		}
	}

	// 等待所有并发任务结束
	wg.Wait()
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
