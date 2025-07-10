package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"gomall/app/seckill/biz/dal/redis"
)

const (
	TaskRollbackStock = "seckill:rollback"
)

func NewRollbackStockTask(userID, activityID string) (*asynq.Task, error) {
	payload := map[string]interface{}{
		"user_id":     userID,
		"activity_id": activityID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskRollbackStock, payloadBytes), nil
}

func HandleRollbackStockTask(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		UserID     string `json:"user_id"`
		ActivityID string `json:"activity_id"`
	}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	userID := payload.UserID
	activityID := payload.ActivityID

	tokenKey := fmt.Sprintf("seckill_token:%s:%s", userID, activityID)
	stockKey := fmt.Sprintf("seckill_stock:%s", activityID)

	// 如果 token 还存在，说明用户未下单，回滚库存
	exist, err := redis.RedisClient.Exists(ctx, tokenKey).Result()
	if err != nil {
		return err
	}
	if exist == 0 {
		// token 已被删除，用户已经结算了
		fmt.Println("已经删除")
		return nil
	}

	// Lua 脚本回滚库存 + 删除 token
	fmt.Println("start to rollback redis stock")
	luaScript := `
		if redis.call("exists", KEYS[1]) == 1 then
			redis.call("incr", KEYS[2])
			redis.call("del", KEYS[1])
			return 1
		end
		return 0
	`

	_, err = redis.RedisClient.Eval(ctx, luaScript, []string{tokenKey, stockKey}).Result()
	if err != nil {
		return err
	}
	return nil
}

func AsyncInit() {
	// 注册任务处理器
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskRollbackStock, HandleRollbackStockTask)

	// 启动异步任务处理器
	go func() {
		if err := AsynqServer.Run(mux); err != nil {
			panic(err)
		}
	}()
}
