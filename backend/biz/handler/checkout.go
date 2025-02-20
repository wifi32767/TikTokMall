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

type CheckoutReq struct {
	UserId  uint32 `json:"user_id" binding:"required"`
	Email   string `json:"email" binding:"required"`
	Address *struct {
		StreetAddress string
		City          string
		State         string
		Country       string
		ZipCode       int32
	} `json:"address"`
	CreditCard *struct {
		CreditCardNumber          string `json:"credit_card_number" binding:"required"`
		CreditCardCvv             int32  `json:"credit_card_cvv" binding:"required"`
		CreditCardExpirationMonth int32  `json:"credit_card_expiration_month" binding:"required"`
		CreditCardExpirationYear  int32  `json:"credit_card_expiration_year" binding:"required"`
	} `json:"credit_card"`
}

type CheckoutResp struct {
	OrderId       string `json:"order_id"`
	TransactionId string `json:"transaction_id"`
}

// @Summary		结算
// @Description	结算购物车，生成订单并支付
// @Tags			Checkout
// @Accept			json
// @Produce		json
// @Param			input	body		CheckoutReq		true	"结算请求"
// @Success		200		{object}	CheckoutResp	"成功"
// @Failure		400		{object}	errorReturn		"请求信息错误"
// @Failure		500		{object}	errorReturn		"服务器错误"
// @Router			/api/checkout [post]
func Checkout(c *gin.Context) {
	input := CheckoutReq{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.CheckoutClient.Checkout(c.Request.Context(), &checkout.CheckoutReq{
		UserId: input.UserId,
		Email:  input.Email,
		Address: &order.Address{
			StreetAddress: input.Address.StreetAddress,
			City:          input.Address.City,
			State:         input.Address.State,
			Country:       input.Address.Country,
			ZipCode:       input.Address.ZipCode,
		},
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          input.CreditCard.CreditCardNumber,
			CreditCardCvv:             input.CreditCard.CreditCardCvv,
			CreditCardExpirationMonth: input.CreditCard.CreditCardExpirationMonth,
			CreditCardExpirationYear:  input.CreditCard.CreditCardExpirationYear,
		},
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}
