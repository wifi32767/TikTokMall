package dal

import (
	"github.com/wifi32767/TikTokMall/app/user/biz/model"
	"github.com/wifi32767/TikTokMall/app/user/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func MysqlInit() {
	dsn := conf.GetConf().Mysql.Dsn
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}
