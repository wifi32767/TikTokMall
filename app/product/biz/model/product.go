package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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

type CachedProductQuery struct {
	productQuery ProductQuery
	client       *redis.Client
}

func NewCachedProductQuery(ctx context.Context, db *gorm.DB, client *redis.Client) *CachedProductQuery {
	return &CachedProductQuery{*NewProductQuery(ctx, db), client}
}

func (c CachedProductQuery) Create(product *Product) error {
	err := c.productQuery.db.WithContext(c.productQuery.ctx).Create(product).Error
	if err == nil {
		cacheKey := fmt.Sprintf("product_%d", product.Model.ID)
		encoded, err := json.Marshal(product)
		if err != nil {
			return err
		}
		_ = c.client.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return err
}

func (c CachedProductQuery) Update(product *Product) error {
	err := c.productQuery.db.WithContext(c.productQuery.ctx).Updates(product).Error
	if err == nil {
		cacheKey := fmt.Sprintf("product_%d", product.Model.ID)
		encoded, err := json.Marshal(product)
		if err != nil {
			return err
		}
		_ = c.client.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return err
}

func (c CachedProductQuery) Delete(id uint32) error {
	err := c.productQuery.db.WithContext(c.productQuery.ctx).Delete(&Product{}, id).Error
	if err == nil {
		cacheKey := fmt.Sprintf("product_%d", id)
		_ = c.client.Del(c.productQuery.ctx, cacheKey)
	}
	return err
}

func (c CachedProductQuery) GetById(id uint32) (*Product, error) {
	cacheKey := fmt.Sprintf("product_%d", id)
	cacheResult := c.client.Get(c.productQuery.ctx, cacheKey)
	product := &Product{}

	err := func() error {
		err1 := cacheResult.Err()
		if err1 != nil {
			return err1
		}
		cachedResultByte, err2 := cacheResult.Bytes()
		if err2 != nil {
			return err2
		}
		err3 := json.Unmarshal(cachedResultByte, &product)
		if err3 != nil {
			return err3
		}
		return nil
	}()
	if err != nil {
		product, err = c.productQuery.GetById(id)
		fmt.Printf("product: %v\n", product)
		if err != nil {
			return nil, err
		}
		encoded, err := json.Marshal(product)
		if err != nil {
			return product, nil
		}
		_ = c.client.Set(c.productQuery.ctx, cacheKey, encoded, time.Hour)
	}
	return product, err
}
