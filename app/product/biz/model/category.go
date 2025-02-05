package model

import (
	"context"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(255);not null;unique"`
	Description string    `gorm:"type:text"`
	Products    []Product `gorm:"many2many:product_category"`
}

func (c Category) TableName() string {
	return "category"
}

type CategoryQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewCategoryQuery(ctx context.Context, db *gorm.DB) *CategoryQuery {
	return &CategoryQuery{ctx, db}
}

// 通过名称获取分类，如果不存在则创建
func (c CategoryQuery) GetByName(name string) (*Category, error) {
	category := &Category{Name: name}
	err := c.db.WithContext(c.ctx).Preload("Products").Model(&Category{}).First(category).Error
	if err == gorm.ErrRecordNotFound {
		err = c.db.WithContext(c.ctx).Create(category).Error
	}
	return category, err
}
