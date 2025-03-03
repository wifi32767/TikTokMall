package dal

import (
	"fmt"

	"github.com/bits-and-blooms/bloom"
	"github.com/wifi32767/TikTokMall/app/product/biz/model"
	"github.com/wifi32767/TikTokMall/app/product/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Filter *bloom.BloomFilter
	err    error
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
	autoMigrate(&model.Product{})
	autoMigrate(&model.Category{})
	Filter = bloom.New(1000000, 5)
	products := make([]model.Product, 0)
	err := DB.Find(&products).Error
	if err != nil {
		panic(err)
	}
	for _, product := range products {
		Filter.AddString(fmt.Sprintf("product_%d", product.Model.ID))
	}
}
