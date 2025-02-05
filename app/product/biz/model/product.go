package model

import (
	"context"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	Picture     string
	Price       float32    `gorm:"type:decimal(10,2)"`
	Categories  []Category `gorm:"many2many:product_category"`
}

func (p Product) TableName() string {
	return "product"
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func NewProductQuery(ctx context.Context, db *gorm.DB) *ProductQuery {
	return &ProductQuery{ctx, db}
}

func (p ProductQuery) Create(product *Product) error {
	err := p.db.WithContext(p.ctx).Create(product).Error
	return err
}

func (p ProductQuery) Update(product *Product) error {
	err := p.db.WithContext(p.ctx).Updates(product).Error
	return err
}

func (p ProductQuery) Delete(id uint32) error {
	err := p.db.WithContext(p.ctx).Delete(&Product{}, id).Error
	return err
}

func (p ProductQuery) GetById(id uint32) (*Product, error) {
	product := &Product{Model: gorm.Model{ID: uint(id)}}
	err := p.db.WithContext(p.ctx).Preload("Categories").Model(&Product{}).First(product).Error
	return product, err
}

func (p ProductQuery) GetAll() (*[]*Product, error) {
	products := &[]*Product{}
	err := p.db.WithContext(p.ctx).Preload("Categories").Find(products).Error
	return products, err
}

func (p ProductQuery) Search(query string) (*[]*Product, error) {
	query = "%" + query + "%"
	products := &[]*Product{}
	err := p.db.WithContext(p.ctx).Preload("Categories").Find(products, "name like ? or description like ?", query, query).Error
	return products, err
}
