package dal

import (
	"github.com/bwmarrin/snowflake"
	"gomall/app/seckill/biz/dal/asynq"
	"gomall/app/seckill/biz/dal/kafka"
	"gomall/app/seckill/biz/dal/mysql"
	"gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/dal/sentinel"
)

var Node *snowflake.Node

func Init() {
	mysql.Init()
	redis.Init()
	sentinel.Init()
	kafka.Init()
	asynq.Init()

	// 雪花算法初始化
	// 设置节点ID（0~1023）
	var err error
	Node, err = snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
}
