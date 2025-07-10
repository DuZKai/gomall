package dal

import (
	"gomall/app/seckill/biz/dal/kafka"
	"gomall/app/seckill/biz/dal/mysql"
	"gomall/app/seckill/biz/dal/redis"
	"gomall/app/seckill/biz/dal/sentinel"
)

var prodNum = 10 // 每秒允许的请求数

func Init() {
	mysql.Init()
	redis.Init()
	sentinel.Init(prodNum)
	kafka.Init()
}
