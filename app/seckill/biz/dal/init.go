package dal

import (
	"gomall/app/seckill/biz/dal/kafka"
	"gomall/app/seckill/biz/dal/sentinel"
)

var prodNum = 10 // 每秒允许的请求数

func Init() {
	// redis.Init()
	// mysql.Init()
	sentinel.Init(prodNum)
	kafka.Init()
}
