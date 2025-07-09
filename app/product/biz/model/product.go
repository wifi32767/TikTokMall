package model

import (
	"context"
	"fmt"
	"time"

	"github.com/bits-and-blooms/bloom"
	"github.com/redis/go-redis/v9"
	"github.com/wifi32767/TikTokMall/common/redis_cache"
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
	cache        *redis_cache.RedisCache
}

func NewCachedProductQuery(ctx context.Context, db *gorm.DB, client *redis.Client, filter *bloom.BloomFilter) *CachedProductQuery {
	return &CachedProductQuery{
		productQuery: *NewProductQuery(ctx, db),
		cache:        redis_cache.NewRedisCache(client, filter),
	}
}

func (c CachedProductQuery) Create(product *Product) error {
	err := c.productQuery.Create(product)
	if err == nil {
		c.cache.Set(c.productQuery.ctx, getKey(uint32(product.Model.ID)), product)
	}
	return err
}

func (c CachedProductQuery) Update(product *Product) error {
	cacheKey := getKey(uint32(product.Model.ID))
	c.cache.Del(c.productQuery.ctx, cacheKey)
	err := c.productQuery.Update(product)
	if err == nil {
		go func() {
			// 延迟双删
			time.Sleep(500 * time.Millisecond)
			c.cache.Del(c.productQuery.ctx, cacheKey)
		}()
	}
	return err
}

func (c CachedProductQuery) Delete(id uint32) error {
	err := c.productQuery.Delete(id)
	if err == nil {
		c.cache.Del(c.productQuery.ctx, getKey(id))
	}
	return err
}

func (c CachedProductQuery) GetById(id uint32) (*Product, error) {
	product := &Product{}
	err := c.cache.Get(c.productQuery.ctx, getKey(id), product)
	if err == redis_cache.CacheMiss {
		// 如果缓存未命中，则从数据库查询
		product, err = c.productQuery.GetById(id)
		if err == nil {
			// 如果数据库查询成功，则将结果存入缓存
			c.cache.Set(c.productQuery.ctx, getKey(id), product)
		}
	}
	return product, err
}

func getKey(id uint32) string {
	return fmt.Sprintf("product_%d", id)
}
