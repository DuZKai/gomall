package dal

import (
	"gomall/app/product/biz/dal/minio"
	"gomall/app/product/biz/dal/mysql"
	"gomall/app/product/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
	minio.Init()
}
