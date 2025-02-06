package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/wifi32767/TikTokMall/app/product/biz/dal"
	"github.com/wifi32767/TikTokMall/app/product/biz/model"
	product "github.com/wifi32767/TikTokMall/rpc/kitex_gen/product"
	"gorm.io/gorm"
)

// ProductCatalogServiceImpl implements the last service interface defined in the IDL.
type ProductCatalogServiceImpl struct{}

// CreateProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) CreateProduct(ctx context.Context, req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	klog.Infof("CreateProduct: %v", req)
	prod := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
	}
	prod.Categories = make([]model.Category, 0)
	// 把所有的category全都拉出来
	query := model.NewCategoryQuery(ctx, dal.DB)
	for _, categoryName := range req.Categories {
		category, err := query.GetByName(categoryName)
		if err != nil {
			klog.Error(err)
			return nil, err
		}
		prod.Categories = append(prod.Categories, *category)
	}
	err = model.NewProductQuery(ctx, dal.DB).Create(prod)
	resp = &product.CreateProductResp{
		Id: uint32(prod.ID),
	}
	return
}

// UpdateProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) UpdateProduct(ctx context.Context, req *product.UpdateProductReq) (resp *product.Empty, err error) {
	klog.Infof("UpdateProduct: %v", req)
	prod := &model.Product{
		Model: gorm.Model{
			ID: uint(req.GetId()),
		},
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Picture:     req.GetPicture(),
		Price:       req.GetPrice(),
	}
	prod.Categories = make([]model.Category, 0)
	query := model.NewCategoryQuery(ctx, dal.DB)
	for _, categoryName := range req.Categories {
		category, err := query.GetByName(categoryName)
		if err != nil {
			klog.Error(err)
			return nil, err
		}
		prod.Categories = append(prod.Categories, *category)
	}
	err = model.NewProductQuery(ctx, dal.DB).Update(prod)
	return
}

// DeleteProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) DeleteProduct(ctx context.Context, req *product.DeleteProductReq) (resp *product.Empty, err error) {
	klog.Infof("DeleteProduct: %v", req)
	err = model.NewProductQuery(ctx, dal.DB).Delete(req.GetId())
	return
}

// ListProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	klog.Infof("ListProducts: %v", req)
	// 如果不设定类别，按照所有类别算
	if req.CategoryName == "" {
		resp = &product.ListProductsResp{
			Products: make([]*product.Product, 0),
		}
		prods, err := model.NewProductQuery(ctx, dal.DB).GetAll()
		if err != nil {
			klog.Error(err)
			return nil, err
		}
		for i := (req.GetPage() - 1) * req.GetPageSize(); i < req.GetPage()*req.GetPageSize() && int(i) < len(*prods); i++ {
			prod := (*prods)[i]
			resp.Products = append(resp.Products, &product.Product{
				Id:          uint32(prod.ID),
				Name:        prod.Name,
				Description: prod.Description,
				Picture:     prod.Picture,
				Price:       prod.Price,
				Categories:  make([]string, 0),
			})
			for _, category := range prod.Categories {
				resp.Products[len(resp.Products)-1].Categories = append(resp.Products[len(resp.Products)-1].Categories, category.Name)
			}
		}
		return resp, err
	}
	category, err := model.NewCategoryQuery(ctx, dal.DB).GetByName(req.GetCategoryName())
	if err != nil {
		klog.Error(err)
		return
	}
	resp = &product.ListProductsResp{
		Products: make([]*product.Product, 0),
	}
	for i := (req.GetPage() - 1) * req.GetPageSize(); i < req.GetPage()*req.GetPageSize() && int(i) < len(category.Products); i++ {
		prod := category.Products[i]
		resp.Products = append(resp.Products, &product.Product{
			Id:          uint32(prod.ID),
			Name:        prod.Name,
			Description: prod.Description,
			Picture:     prod.Picture,
			Price:       prod.Price,
			Categories:  make([]string, 0),
		})
		for _, category := range prod.Categories {
			resp.Products[len(resp.Products)-1].Categories = append(resp.Products[len(resp.Products)-1].Categories, category.Name)
		}
	}
	return
}

// GetProduct implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	klog.Infof("GetProduct: %v", req)
	prod, err := model.NewCachedProductQuery(ctx, dal.DB, dal.RedisClient).GetById(req.GetId())
	if err != nil {
		klog.Error(err)
		return
	}
	resp = &product.GetProductResp{
		Product: &product.Product{
			Id:          uint32(prod.ID),
			Name:        prod.Name,
			Description: prod.Description,
			Picture:     prod.Picture,
			Price:       prod.Price,
			Categories:  make([]string, 0),
		},
	}
	for _, category := range prod.Categories {
		resp.Product.Categories = append(resp.Product.Categories, category.Name)
	}
	return
}

// SearchProducts implements the ProductCatalogServiceImpl interface.
func (s *ProductCatalogServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	klog.Infof("SearchProducts: %v", req)
	prods, err := model.NewProductQuery(ctx, dal.DB).Search(req.Query)
	if err != nil {
		klog.Error(err)
		return
	}
	var results []*product.Product
	for _, prod := range *prods {
		results = append(results, &product.Product{
			Id:          uint32(prod.ID),
			Name:        prod.Name,
			Description: prod.Description,
			Picture:     prod.Picture,
			Price:       prod.Price,
			Categories:  make([]string, 0),
		})
		for _, category := range prod.Categories {
			results[len(results)-1].Categories = append(results[len(results)-1].Categories, category.Name)
		}
	}
	resp = &product.SearchProductsResp{
		Results: results,
	}
	return
}
