package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gomall/app/product/conf"
)

var (
	RedisClient      *redis.Client
	RedisBloomClient *redis.Client
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

	RedisBloomClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().RedisBloom.Address,
		Username: conf.GetConf().RedisBloom.Username,
		Password: conf.GetConf().RedisBloom.Password,
		DB:       conf.GetConf().RedisBloom.DB,
	})
	if err := RedisBloomClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
