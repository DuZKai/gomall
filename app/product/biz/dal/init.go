package dal

import (
	"gomall/app/product/biz/dal/sentinel"
	"gomall/app/product/rpc"
)

func Init() {
	// redis.Init()
	// mysql.Init()
	// minio.Init()
	// es.Init()
	// oss.Init()
	sentinel.Init()
	rpc.InitUserRPC()
}
