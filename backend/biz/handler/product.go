package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/product"
)

// @Summary		创建商品
// @Description	创建一个新的商品
// @Tags			product
// @Produce		json
// @Param			req	body		productCreateReq	true	"商品信息"
// @Success		200		{object}	idReq				"商品id
// @Failure		400		{object}	errorReturn			"请求格式错误"
// @Failure		500		{object}	errorReturn			"服务器错误"
// @Router			/product/create [post]
func ProductCreate(c *gin.Context) {
	req := productCreateReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.ProductClient.CreateProduct(c.Request.Context(), &product.CreateProductReq{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  req.Categories,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": resp.GetId(),
	})
}

// @Summary		更新商品信息
// @Description	更新商品信息
// @Tags			product
// @Produce		json
// @Param			req	body	simpleProduct	true	"商品信息"
// @Success		200		"成功"
// @Failure		400		{object}	errorReturn	"请求格式错误"
// @Failure		500		{object}	errorReturn	"服务器错误"
// @Router			/product/update [put]
func ProductUpdate(c *gin.Context) {
	req := simpleProduct{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.ProductClient.UpdateProduct(c.Request.Context(), &product.UpdateProductReq{
		Id:          req.ID,
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  req.Categories,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// @Summary		删除商品
// @Description	删除商品
// @Tags			product
// @Produce		json
// @Param			req	body	idReq	true	"商品id"
// @Success		200		"成功"
// @Failure		400		{object}	errorReturn	"请求格式错误"
// @Failure		500		{object}	errorReturn	"服务器错误"
// @Router			/product/delete [delete]
func ProductDelete(c *gin.Context) {
	req := idReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.ProductClient.DeleteProduct(c.Request.Context(), &product.DeleteProductReq{
		Id: req.ID,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// @Summary		商品列表
// @Description	获取指定类别的商品列表，不指定类别返回全部商品的列表
// @Tags			product
// @Produce		json
// @Param			req	body		productListReq	true	"商品列表信息"
// @Success		200		{object}	productListReq	"商品列表"
// @Failure		400		{object}	errorReturn			"请求格式错误"
// @Failure		500		{object}	errorReturn			"服务器错误"
// @Router			/product/list [get]
func ProductList(c *gin.Context) {
	req := productListReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.ProductClient.ListProducts(c.Request.Context(), &product.ListProductsReq{
		PageSize:     req.PageSize,
		Page:         req.Page,
		CategoryName: req.CategoryName,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	products := make([]simpleProduct, 0)
	for _, prod := range resp.GetProducts() {
		products = append(products, simpleProduct{
			ID:          prod.GetId(),
			Name:        prod.GetName(),
			Description: prod.GetDescription(),
			Picture:     prod.GetPicture(),
			Price:       prod.GetPrice(),
			Categories:  prod.GetCategories(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// @Summary		获取商品信息
// @Description	获取单个商品信息
// @Tags			product
// @Produce		json
// @Param			req	body		idReq			true	"商品id"
// @Success		200		{object}	simpleProduct	"商品信息"
// @Failure		400		{object}	errorReturn		"请求格式错误"
// @Failure		500		{object}	errorReturn		"服务器错误"
// @Router			/product/get [get]
func ProductGet(c *gin.Context) {
	req := idReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.ProductClient.GetProduct(c.Request.Context(), &product.GetProductReq{
		Id: req.ID,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	prod := resp.GetProduct()
	c.JSON(http.StatusOK, gin.H{
		"product": simpleProduct{
			ID:          prod.GetId(),
			Name:        prod.GetName(),
			Description: prod.GetDescription(),
			Picture:     prod.GetPicture(),
			Price:       prod.GetPrice(),
			Categories:  prod.GetCategories(),
		},
	})
}

// @Summary		商品搜索
// @Description	搜索在名字和描述中含有关键词的商品
// @Tags			product
// @Produce		json
// @Param			req	body		searchReq		true	"搜索关键词"
// @Success		200		{object}	simpleProduct	"商品信息"
// @Failure		400		{object}	errorReturn		"请求格式错误"
// @Failure		500		{object}	errorReturn		"服务器错误"
// @Router			/product/search [get]
func ProductSearch(c *gin.Context) {
	req := searchReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.ProductClient.SearchProducts(c.Request.Context(), &product.SearchProductsReq{
		Query: req.Query,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	results := make([]simpleProduct, 0)
	for _, prod := range resp.GetResults() {
		results = append(results, simpleProduct{
			ID:          prod.GetId(),
			Name:        prod.GetName(),
			Description: prod.GetDescription(),
			Picture:     prod.GetPicture(),
			Price:       prod.GetPrice(),
			Categories:  prod.GetCategories(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"results": results,
	})
}
