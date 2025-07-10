package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"log"
	"time"
)

// 实现 sarama.ConsumerGroupHandler 接口
type SeckillConsumer struct{}

func (h *SeckillConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *SeckillConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *SeckillConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.Background()
	for msg := range claim.Messages() {
		var req model.SeckillMessage
		err := json.Unmarshal(msg.Value, &req)
		if err != nil {
			log.Printf("[Consumer] Unmarshal error: %v", err)
			continue
		}

		log.Printf("[Consumer] Processing user %s in activity %s at %s", req.UserID, req.ActivityID, req.TS)

		// 黑名单校验
		blacklistKey := fmt.Sprintf("blacklist:%s", req.UserID)
		isBlacklisted, err := redis.RedisClient.Exists(ctx, blacklistKey).Result()
		if err != nil {
			log.Printf("[Consumer] Redis error during blacklist check: %v", err)
			continue
		}
		if isBlacklisted > 0 {
			log.Printf("[Consumer] User %s is blacklisted, dropping", req.UserID)
			redis.RedisClient.Set(ctx, fmt.Sprintf("seckill_fail:%s:%s", req.UserID, req.ActivityID), 1, time.Minute)
			sess.MarkMessage(msg, "")
			continue
		}

		// 访问频率校验
		freqKey := fmt.Sprintf("freq:%s", req.UserID)
		count, err := redis.RedisClient.Incr(ctx, freqKey).Result()
		if err != nil {
			log.Printf("[Consumer] Redis error during freq incr: %v", err)
			continue
		}
		if count == 1 {
			redis.RedisClient.Expire(ctx, freqKey, 10*time.Second)
		}
		if count > 5 {
			log.Printf("[Consumer] User %s is too frequent (%d), adding to blacklist", req.UserID, count)
			redis.RedisClient.Set(ctx, blacklistKey, 1, 30*time.Minute) // 拉黑 30 分钟
			redis.RedisClient.Set(ctx, fmt.Sprintf("seckill_fail:%s:%s", req.UserID, req.ActivityID), 1, time.Minute)
			sess.MarkMessage(msg, "")
			continue
		}

		// 第六步：库存判断，生成 token
		luaScript := `
			if redis.call("get", KEYS[2]) then
				return {1, redis.call("get", KEYS[2])}
			end
			local stock = tonumber(redis.call("get", KEYS[1]))
			if not stock or stock <= 0 then
				return {0, "stock_empty"}
			end
			redis.call("decr", KEYS[1])
			redis.call("setex", KEYS[2], tonumber(ARGV[2]), ARGV[1])
			return {2, ARGV[1]}
		`

		// 准备 key 和参数
		stockKey := fmt.Sprintf("seckill_stock:%s", req.ActivityID)
		tokenKey := fmt.Sprintf("seckill_token:%s:%s", req.UserID, req.ActivityID)
		tokenVal := fmt.Sprintf("token:%s:%s:%d", req.UserID, req.ActivityID, time.Now().UnixNano())
		tokenTTL := 30 * 60 // 30分钟
		result, err := redis.RedisClient.Eval(ctx, luaScript, []string{stockKey, tokenKey}, tokenVal, tokenTTL).Result()
		if err != nil {
			log.Printf("[Consumer] Lua eval failed: %v", err)
			continue
		}

		resArr, ok := result.([]interface{})
		if !ok || len(resArr) < 2 {
			log.Printf("[Consumer] Lua return parse failed")
			continue
		}

		statusCode := int(resArr[0].(int64))
		returnToken := resArr[1].(string)

		switch statusCode {
		case 0:
			log.Printf("[Consumer] Stock empty for activity %s", req.ActivityID)
			redis.RedisClient.Set(ctx, fmt.Sprintf("seckill_fail:%s:%s", req.UserID, req.ActivityID), 1, time.Minute)
		case 1:
			log.Printf("[Consumer] Duplicate token for user %s: %s", req.UserID, returnToken)
		case 2:
			log.Printf("[Consumer] Success, token for user %s: %s", req.UserID, returnToken)
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}
