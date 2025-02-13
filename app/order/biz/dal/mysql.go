package dal

import (
	"github.com/wifi32767/TikTokMall/app/order/biz/model"
	"github.com/wifi32767/TikTokMall/app/order/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func autoMigrate(model any) {
	err = DB.AutoMigrate(model)
	if err != nil {
		panic(err)
	}
}

func MysqlInit() {
	dsn := conf.GetConf().Mysql.Dsn
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	autoMigrate(&model.Order{})
	autoMigrate(&model.OrderItem{})
}
