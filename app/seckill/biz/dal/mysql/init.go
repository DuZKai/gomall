package mysql

import (
	"fmt"
	"gomall/app/seckill/conf"
	"gomall/app/user/biz/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))
	fmt.Println("DSN:", dsn)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	if conf.GetConf().Env != "online" {
		err = DB.AutoMigrate(&model.User{})
		if err != nil {
			panic(err)
		}
	}
}
