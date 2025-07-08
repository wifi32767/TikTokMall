package model

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bits-and-blooms/bloom"
	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string     `gorm:"type:varchar(255);not null"`
	Description string     `gorm:"type:text"`
	Picture     string     `gorm:"type:mediumtext"`
	Price       float32    `gorm:"type:decimal(10,2)"`
	Categories  []Category `gorm:"many2many:product_category"`
}

func (p Product) TableName() string {
	return "products"
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
	res := p.db.WithContext(p.ctx).Where("id = ?", product.ID).Updates(product)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

func (p ProductQuery) Delete(id uint32) error {
	res := p.db.WithContext(p.ctx).Delete(&Product{}, id)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}

func (p ProductQuery) GetById(id uint32) (*Product, error) {
	product := &Product{}
	err := p.db.WithContext(p.ctx).Preload("Categories").Model(&Product{}).Where("id = ?", id).First(product).Error
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
	filter       *bloom.BloomFilter
}

func NewCachedProductQuery(ctx context.Context, db *gorm.DB, client *redis.Client, filter *bloom.BloomFilter) *CachedProductQuery {
	return &CachedProductQuery{*NewProductQuery(ctx, db), client, filter}
}

func (c CachedProductQuery) Create(product *Product) error {
	err := c.productQuery.Create(product)
	if err == nil {
		cacheKey := fmt.Sprintf("product_%d", product.Model.ID)
		c.filter.AddString(cacheKey)
		encoded, err := json.Marshal(product)
		if err != nil {
			return err
		}
		_ = c.client.Set(c.productQuery.ctx, cacheKey, encoded, randomDuration())
	}

	return err
}

func (c CachedProductQuery) Update(product *Product) error {
	cacheKey := fmt.Sprintf("product_%d", product.Model.ID)
	_ = c.client.Del(c.productQuery.ctx, cacheKey)
	err := c.productQuery.Update(product)
	if err == nil {
		go func() {
			// 延迟双删
			time.Sleep(500 * time.Millisecond)
			_ = c.client.Del(c.productQuery.ctx, cacheKey)
		}()
	}
	return err
}

func (c CachedProductQuery) Delete(id uint32) error {
	err := c.productQuery.Delete(id)
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
		// 续期
		_ = c.client.Expire(c.productQuery.ctx, cacheKey, randomDuration())
		err3 := json.Unmarshal(cachedResultByte, &product)
		if err3 != nil {
			return err3
		}
		return nil
	}()
	// 缓存未命中，或者缓存解析失败
	if err != nil {
		if !c.filter.TestString(cacheKey) {
			return nil, gorm.ErrRecordNotFound
		}
		product, err = c.productQuery.GetById(id)
		if err != nil {
			return nil, err
		}
		encoded, err := json.Marshal(product)
		if err != nil {
			return product, nil
		}
		_ = c.client.Set(c.productQuery.ctx, cacheKey, encoded, randomDuration())
	}
	return product, err
}

// 在失效时间内加入随机数，可以防止缓存雪崩的情况
func randomDuration() time.Duration {
	return time.Duration(1+rand.Intn(2)) * time.Hour
}
