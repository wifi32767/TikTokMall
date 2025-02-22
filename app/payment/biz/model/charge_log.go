package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type PaymentLog struct {
	gorm.Model
	UserId        uint32
	OrderId       string
	TransactionId string
	Amount        float32
	PayAt         time.Time
}

func (p PaymentLog) TableName() string {
	return "payment_logs"
}

func CreatePaymentLog(db *gorm.DB, ctx context.Context, p *PaymentLog) error {
	return db.WithContext(ctx).Create(p).Error
}
