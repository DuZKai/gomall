package dal

import (
	"gomall/app/user/biz/dal/etcd"
	"gomall/app/user/biz/dal/mysql"
)

func Init() {
	mysql.Init()
	etcd.Init()
}
