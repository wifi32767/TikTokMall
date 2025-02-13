package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart"
)

type cartItem struct {
	ProductId uint32 `json:"product_id" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required"`
}

type addItemReq struct {
	UserId uint32   `json:"userid" binding:"required"`
	Item   cartItem `json:"item" binding:"required"`
}

type userIdReq struct {
	UserId uint32 `json:"userid" binding:"required"`
}

//	@Summary		添加商品到购物车
//	@Description	添加商品到指定用户的购物车
//	@Tags			cart
//	@Produce		json
//	@Param			input	body	addItemReq	true	"商品和用户信息"
//	@Success		200		"成功"
//	@Failure		400		{object}	errorReturn	"请求格式错误"
//	@Failure		500		{object}	errorReturn	"服务器错误"
//	@Router			/cart/additem [post]
func CartAddItem(c *gin.Context) {
	req := addItemReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.CartClient.AddItem(c.Request.Context(), &cart.AddItemReq{
		UserId: req.UserId,
		Item: &cart.CartItem{
			ProductId: req.Item.ProductId,
			Quantity:  req.Item.Quantity,
		},
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

//	@Summary		获取购物车
//	@Description	获取指定用户的购物车中的商品
//	@Tags			cart
//	@Produce		json
//	@Param			input	body		userIdReq	true	"用户id"
//	@Success		200		{object}	[]cartItem	"购物车中的商品"
//	@Failure		400		{object}	errorReturn	"请求格式错误"
//	@Failure		500		{object}	errorReturn	"服务器错误"
//	@Router			/cart/get [get]
func CartGet(c *gin.Context) {
	req := userIdReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.CartClient.GetCart(c.Request.Context(), &cart.GetCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": resp.Cart.Items,
	})
}

//	@Summary		清空购物车
//	@Description	清空指定用户的购物车
//	@Tags			cart
//	@Produce		json
//	@Param			input	body	userIdReq	true	"用户id"
//	@Success		200		"成功"
//	@Failure		400		{object}	errorReturn	"请求格式错误"
//	@Failure		500		{object}	errorReturn	"服务器错误"
//	@Router			/cart/empty [post]
func CartEmpty(c *gin.Context) {
	req := userIdReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.CartClient.EmptyCart(c.Request.Context(), &cart.EmptyCartReq{
		UserId: req.UserId,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
