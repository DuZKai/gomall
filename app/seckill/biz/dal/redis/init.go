package redis

import (
	"context"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"gomall/app/seckill/conf"
)

var (
	RedisClient *redis.Client
	Locker      *redislock.Client
)

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Username: conf.GetConf().Redis.Username,
		Password: conf.GetConf().Redis.Password,
		DB:       conf.GetConf().Redis.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	Locker = redislock.New(RedisClient)
}
