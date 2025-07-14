package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/model"
	"gomall/app/seckill/config"
	"log"
	"time"
)

// 实现 sarama.ConsumerGroupHandler 接口
type SeckillConsumer struct{}

func (h *SeckillConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *SeckillConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *SeckillConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := context.Background()

	// 整合所有校验和秒杀逻辑的 Lua 脚本
	megaScript := `
		-- 1. 黑名单检查
		if redis.call("EXISTS", KEYS[1]) > 0 then
			return {1, "blacklisted"}  -- 状态码1: 黑名单
		end
		
		-- 2. 频率检查
		local freq = redis.call("INCR", KEYS[2])
		if freq == 1 then
			redis.call("EXPIRE", KEYS[2], ARGV[1])
		end
		if freq > tonumber(ARGV[2]) then
			redis.call("SETEX", KEYS[1], ARGV[3], "1")  -- 加入黑名单
			return {2, "freq_limit"}  -- 状态码2: 频率限制
		end
		
		-- 3. 幂等检查
		local idemp = redis.call("SETNX", KEYS[3], "1")
		if idemp == 0 then
			return {3, "duplicate"}  -- 状态码3: 重复请求
		end
		redis.call("EXPIRE", KEYS[3], ARGV[4])
		
		-- 4. Token 存在检查
		if redis.call("EXISTS", KEYS[5]) > 0 then
			return {4, redis.call("GET", KEYS[5])}  -- 状态码4: 已有token
		end
		
		-- 5. 库存检查
		local stock = tonumber(redis.call("GET", KEYS[4]))
		if not stock or stock <= 0 then
			return {0, "stock_empty"}  -- 状态码0: 库存不足
		end
		
		-- 6. 扣减库存并生成 token
		redis.call("DECR", KEYS[4])
		redis.call("SETEX", KEYS[5], ARGV[6], ARGV[5])
		return {5, ARGV[5]}  -- 状态码5: 成功生成token
	`

	scriptSHA, err := redis.RedisClient.ScriptLoad(ctx, megaScript).Result()
	if err != nil {
		log.Printf("[Consumer] Failed to load mega Lua script: %v", err)
	}

	for msg := range claim.Messages() {
		var req model.SeckillMessage
		if err := json.Unmarshal(msg.Value, &req); err != nil {
			log.Printf("[Consumer] Unmarshal error: %v", err)
			continue
		}

		log.Printf("[Consumer] Processing user %s in activity %s", req.UserID, req.ActivityID)

		// 准备所有 keys 和 args
		keys := []string{
			fmt.Sprintf("seckill:blacklist:%s", req.UserID),                      // KEYS[1]
			fmt.Sprintf("freq:%s", req.UserID),                                   // KEYS[2]
			fmt.Sprintf("seckill:msg:%s:%s", req.UserID, req.ActivityID),         // KEYS[3]
			fmt.Sprintf("seckill:stock:%s", req.ActivityID),                      // KEYS[4]
			fmt.Sprintf("seckill:token:valid:%s:%s", req.ActivityID, req.UserID), // KEYS[5]
		}

		token := model.TokenInfo{
			UserID:       req.UserID,
			ActivityID:   req.ActivityID,
			CreateTime:   time.Now().UnixNano(),
			ExpireSecond: int64(config.AppConfig.TokenTTL * 60),
		}
		tokenData, _ := json.Marshal(token)

		args := []interface{}{
			config.AppConfig.FreqLimitExpire,   // ARGV[1]: 频率过期时间(秒)
			5,                                  // ARGV[2]: 频率阈值
			config.AppConfig.BlacklistTTL * 60, // ARGV[3]: 黑名单TTL(秒)
			config.AppConfig.IdempotentKeyExpire * 60, // ARGV[4]: 幂等key TTL(秒)
			tokenData,                          // ARGV[5]: token 数据
			config.AppConfig.TokenTTL * 2 * 60, // ARGV[6]: token TTL(秒)
		}

		// 执行整合脚本
		res, err := redis.RedisClient.EvalSha(ctx, scriptSHA, keys, args...).Result()
		if err != nil {
			res, err = redis.RedisClient.Eval(ctx, megaScript, keys, args...).Result()
			if err != nil {
				log.Printf("[Consumer] Mega Lua eval failed: %v", err)
				continue
			}
		}

		// 解析结果
		arr, ok := res.([]interface{})
		if !ok || len(arr) < 2 {
			log.Printf("[Consumer] Unexpected mega Lua result: %#v", res)
			continue
		}

		status := int(arr[0].(int64))
		msgData := arr[1].(string)

		// 处理不同状态
		switch status {
		case 0: // 库存不足
			log.Printf("[Consumer] Stock empty for activity %s", req.ActivityID)
			redis.RedisClient.Set(ctx, fmt.Sprintf("seckill:fail:%s:%s", req.ActivityID, req.UserID), 1, time.Minute)
		case 1: // 黑名单
			log.Printf("[Consumer] User %s is blacklisted", req.UserID)
		case 2: // 频率超限
			log.Printf("[Consumer] User %s exceeded frequency limit", req.UserID)
		case 3: // 重复请求
			log.Printf("[Consumer] Duplicate request for user %s", req.UserID)
		case 4: // 已有token
			log.Printf("[Consumer] Existing token for user %s: %s", req.UserID, msgData)
		case 5: // 成功
			log.Printf("[Consumer] Success, token for user %s: %s", req.UserID, msgData)
		}

		// 标记消息已处理
		sess.MarkMessage(msg, "")
	}
	return nil
}
