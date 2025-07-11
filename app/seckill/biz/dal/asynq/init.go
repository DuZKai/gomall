package asynq

import (
	"fmt"
	"github.com/hibiken/asynq"
	"gomall/app/seckill/conf"
)

var (
	AsynqClient    *asynq.Client
	AsynqServer    *asynq.Server
	AsynqScheduler *asynq.Scheduler
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

	// 启动定时任务调度器
	AsynqScheduler = asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{})
	_, err := AsynqScheduler.Register("@every 10s", NewRollbackSchedulerTask())
	if err != nil {
		panic(fmt.Sprintf("Failed to register rollback scheduler task: %v", err))
	}
	go func() {
		if err := AsynqScheduler.Run(); err != nil {
			panic(err)
		}
	}()
	AsyncInit()
}
