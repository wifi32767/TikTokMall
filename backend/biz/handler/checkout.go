package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/checkout"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/payment"
)

// @Summary		结算
// @Description	结算购物车，生成订单并支付
// @Tags			Checkout
// @Accept			json
// @Produce		json
// @Param			req	body		checkoutReq		true	"结算请求"
// @Success		200		{object}	checkoutResp	"成功"
// @Failure		400		{object}	errorReturn		"请求信息错误"
// @Failure		500		{object}	errorReturn		"服务器错误"
// @Router			/api/checkout [post]
func Checkout(c *gin.Context) {
	req := checkoutReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userid := c.GetUint("userid")
	if req.UserId != uint32(userid) {
		permission := c.GetInt("permission")
		if permission != 2 {
			c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
			return
		}
	}
	resp, err := rpc.CheckoutClient.Checkout(c.Request.Context(), &checkout.CheckoutReq{
		UserId: req.UserId,
		Email:  req.Email,
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
		},
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
		},
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
