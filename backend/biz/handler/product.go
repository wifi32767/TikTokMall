package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/product"
)

type productCreateInput struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Picture     string   `json:"picture" binding:"required"`
	Price       float32  `json:"price" binding:"required"`
	Categories  []string `json:"categories"`
}

type simpleProduct struct {
	ID          uint32   `json:"id" binding:"required"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	Price       float32  `json:"price"`
	Categories  []string `json:"categories"`
}

type idInput struct {
	ID uint32 `json:"id" binding:"required"`
}

type productListInput struct {
	PageSize     int32  `json:"pagesize" binding:"required"`
	Page         int32  `json:"page" binding:"required"`
	CategoryName string `json:"category_name"`
}

type searchInput struct {
	Query string `json:"query" binding:"required"`
}

// @Summary 创建商品
// @Description 创建一个新的商品
// @Tags product
// @Produce json
// @Param input body productCreateInput true "商品信息"
// @Success 200 {object} idInput "商品id
// @Failure 400 {object} errorReturn "请求格式错误"
// @Failure 500 {object} errorReturn "服务器错误"
// @Router /product/create [post]
func ProductCreate(c *gin.Context) {
	input := productCreateInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.ProductClient.CreateProduct(c.Request.Context(), &product.CreateProductReq{
		Name:        input.Name,
		Description: input.Description,
		Picture:     input.Picture,
		Price:       input.Price,
		Categories:  input.Categories,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": resp.GetId(),
	})
}

// @Summary 更新商品信息
// @Description 更新商品信息
// @Tags product
// @Produce json
// @Param input body simpleProduct true "商品信息"
// @Success 200 "成功"
// @Failure 400 {object} errorReturn "请求格式错误"
// @Failure 500 {object} errorReturn "服务器错误"
// @Router /product/update [put]
func ProductUpdate(c *gin.Context) {
	input := simpleProduct{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.ProductClient.UpdateProduct(c.Request.Context(), &product.UpdateProductReq{
		Id:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Picture:     input.Picture,
		Price:       input.Price,
		Categories:  input.Categories,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// @Summary 删除商品
// @Description 删除商品
// @Tags product
// @Produce json
// @Param input body idInput true "商品id"
// @Success 200 "成功"
// @Failure 400 {object} errorReturn "请求格式错误"
// @Failure 500 {object} errorReturn "服务器错误"
// @Router /product/delete [delete]
func ProductDelete(c *gin.Context) {
	input := idInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.ProductClient.DeleteProduct(c.Request.Context(), &product.DeleteProductReq{
		Id: input.ID,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// @Summary 商品列表
// @Description 获取指定类别的商品列表，不指定类别返回全部商品的列表
// @Tags product
// @Produce json
// @Param input body productListInput true "商品列表信息"
// @Success 200 {object} productListInput "商品列表"
// @Failure 400 {object} errorReturn "请求格式错误"
// @Failure 500 {object} errorReturn "服务器错误"
// @Router /product/list [get]
func ProductList(c *gin.Context) {
	input := productListInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.ProductClient.ListProducts(c.Request.Context(), &product.ListProductsReq{
		PageSize:     input.PageSize,
		Page:         input.Page,
		CategoryName: input.CategoryName,
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

// @Summary 获取商品信息
// @Description 获取单个商品信息
// @Tags product
// @Produce json
// @Param input body idInput true "商品id"
// @Success 200 {object} simpleProduct "商品信息"
// @Failure 400 {object} errorReturn "请求格式错误"
// @Failure 500 {object} errorReturn "服务器错误"
// @Router /product/get [get]
func ProductGet(c *gin.Context) {
	input := idInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.ProductClient.GetProduct(c.Request.Context(), &product.GetProductReq{
		Id: input.ID,
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

// @Summary 商品搜索
// @Description 搜索在名字和描述中含有关键词的商品
// @Tags product
// @Produce json
// @Param input body searchInput true "搜索关键词"
// @Success 200 {object} simpleProduct "商品信息"
// @Failure 400 {object} errorReturn "请求格式错误"
// @Failure 500 {object} errorReturn "服务器错误"
// @Router /product/search [get]
func ProductSearch(c *gin.Context) {
	input := searchInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.ProductClient.SearchProducts(c.Request.Context(), &product.SearchProductsReq{
		Query: input.Query,
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
