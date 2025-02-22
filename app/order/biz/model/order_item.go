package model

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	OrderIdRefer string `gorm:"size:256;index"`
	ProductId    uint32
	Quantity     int32
	Cost         float32
}

func (o *OrderItem) TableName() string {
	return "order_items"
}
