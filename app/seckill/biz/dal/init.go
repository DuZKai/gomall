package dal

import (
	"gomall/app/seckill/biz/dal/sentinel"
)

func Init() {
	// redis.Init()
	// mysql.Init()
	sentinel.Init()
}
