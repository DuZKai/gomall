package asynq

import (
	"github.com/hibiken/asynq"
	"gomall/app/seckill/conf"
)

var (
	AsynqClient *asynq.Client
	AsynqServer *asynq.Server
)

func Init() {
	redisOpt := asynq.RedisClientOpt{Addr: conf.GetConf().Redis.Address}

	AsynqClient = asynq.NewClient(redisOpt)

	AsynqServer = asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})
	AsyncInit()
}
