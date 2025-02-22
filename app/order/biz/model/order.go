package model

import (
	"context"

	"gorm.io/gorm"
)

type Consignee struct {
	Email string `gorm:"type:varchar(255)"`

	StreetAddress string `gorm:"type:varchar(255)"`
	City          string `gorm:"type:varchar(255)"`
	State         string `gorm:"type:varchar(255)"`
	Country       string `gorm:"type:varchar(255)"`
	ZipCode       int32
}

type OrderState string

const (
	OrderStatePlaced   OrderState = "placed"
	OrderStatePaid     OrderState = "paid"
	OrderStateCanceled OrderState = "canceled"
	// 其实超时自动取消和普通的取消应该是不同的状态
)

type Order struct {
	gorm.Model
	OrderId      string `gorm:"uniqueIndex;type:varchar(255)"`
	UserId       uint32
	UserCurrency string      `gorm:"type:varchar(5)"`
	Consignee    Consignee   `gorm:"embedded"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"`
	OrderState   OrderState  `gorm:"type:varchar(16)"`
}

func (o *Order) TableName() string {
	return "orders"
}

func PlaceOrder(db *gorm.DB, ctx context.Context, order *Order, orderItemList *[]OrderItem) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if err := tx.Create(orderItemList).Error; err != nil {
			return err
		}
		return nil
	})
}

func ListOrder(db *gorm.DB, ctx context.Context, userid uint32) (*[]Order, error) {
	orders := []Order{}
	err := db.WithContext(ctx).Preload("OrderItems").Where("user_id = ?", userid).Find(&orders).Error
	return &orders, err
}

func GetOrder(db *gorm.DB, ctx context.Context, order_id string) (Order, error) {
	order := Order{}
	err := db.WithContext(ctx).Where("order_id = ?", order_id).First(&order).Error
	return order, err
}

func UpdateOrderState(db *gorm.DB, ctx context.Context, userid uint32, order_id string, state OrderState) error {
	return db.WithContext(ctx).Model(&Order{}).Where("order_id = ? AND user_id = ?", order_id, userid).Update("order_state", state).Error
}
