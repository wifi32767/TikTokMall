package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order"
)

const (
	OrderStatePlaced   uint32 = 0
	OrderStatePaid     uint32 = 1
	OrderStateCanceled uint32 = 2
)

type OrderPlaceReq struct {
	UserId       uint32 `json:"user_id" binding:"required"`
	UserCurrency string `json:"user_currency" binding:"required"`
	Address      struct {
		StreetAddress string `json:"street_address" binding:"required"`
		City          string `json:"city" binding:"required"`
		State         string `json:"state" binding:"required"`
		Country       string `json:"country" binding:"required"`
		ZipCode       int32  `json:"zip_code" binding:"required"`
	} `json:"address"`
	Items []struct {
		ProductId uint32  `json:"product_id" binding:"required"`
		Quantity  int32   `json:"quantity" binding:"required"`
		Cost      float32 `json:"cost" binding:"required"`
	} `json:"items" binding:"required"`
}

type UserIdReq struct {
	UserId uint32 `json:"user_id" binding:"required"`
}

type OrderIdReq struct {
	UserId  uint32 `json:"user_id" binding:"required"`
	OrderId string `json:"order_id" binding:"required"`
}

type OrderListRes struct {
	Orders []*struct {
		OrderId      string
		UserId       uint32
		UserCurrency string
		Address      *struct {
			StreetAddress string
			City          string
			State         string
			Country       string
			ZipCode       int32
		}
		Email     string
		CreatedAt int32
	}
}

// @Summary		创建订单
// @Description	创建一个新的订单
// @Tags			order
// @Produce		json
// @Param			input	body		OrderPlaceReq	true	"订单信息"
// @Success		200		{object}	OrderIdReq		"订单id"
// @Failure		400		{object}	errorReturn		"请求格式错误"
// @Failure		500		{object}	errorReturn		"服务器错误"
// @Router			/order/place [post]
func OrderPlace(c *gin.Context) {
	req := OrderPlaceReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items := make([]*order.OrderItem, 0)
	for _, item := range req.Items {
		items = append(items, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			},
			Cost: item.Cost,
		})
	}
	resp, err := rpc.OrderClient.PlaceOrder(c.Request.Context(), &order.PlaceOrderReq{
		UserId:       req.UserId,
		UserCurrency: req.UserCurrency,
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
		},
		OrderItems: items,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"order_id": resp.Order.GetOrderId()})
}

// @Summary		订单列表
// @Description	获取用户的订单列表
// @Tags			order
// @Produce		json
// @Param			input	body		UserIdReq		true	"用户id"
// @Success		200		{object}	OrderListRes	"订单列表"
// @Failure		400		{object}	errorReturn		"请求格式错误"
// @Failure		500		{object}	errorReturn		"服务器错误"
// @Router			/order/list [get]
func OrderList(c *gin.Context) {
	req := UserIdReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.OrderClient.ListOrder(c.Request.Context(), &order.ListOrderReq{
		UserId: req.UserId,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	// 这里注意时间是int32，需要前端转换成time
	c.JSON(http.StatusOK, gin.H{"orders": resp.Orders})
}

// @Summary:		取消订单
// @Description:	取消订单
// @Tags			order
// @Produce		json
// @Param			input	body	OrderIdReq	true	"订单id"
// @Success		200		"成功"
// @Failure		400		{object}	errorReturn	"请求格式错误"
// @Failure		500		{object}	errorReturn	"服务器错误"
// @Router			/order/cancel [put]
func OrderCancel(c *gin.Context) {
	req := OrderIdReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.OrderClient.UpdateOrderState(c.Request.Context(), &order.UpdateOrderStateReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
		State:   OrderStateCanceled,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.Status(http.StatusOK)
}
