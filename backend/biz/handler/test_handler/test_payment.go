package test_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/payment"
)

type PaymentChargeReq struct {
	UserId     uint32     `json:"user_id" binding:"required"`
	OrderId    string     `json:"order_id" binding:"required"`
	Amount     float32    `json:"amount" binding:"required"`
	CreditCard CreditCard `json:"credit_card" binding:"required"`
}

type CreditCard struct {
	CreditCardNumber          string `json:"credit_card_number" binding:"required"`
	CreditCardCvv             int32  `json:"credit_card_cvv" binding:"required"`
	CreditCardExpirationMonth int32  `json:"credit_card_expiration_month" binding:"required"`
	CreditCardExpirationYear  int32  `json:"credit_card_expiration_year" binding:"required"`
}

func PaymentCharge(c *gin.Context) {
	var req PaymentChargeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}
	_, err := rpc.PaymentClient.Charge(c, &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
		Amount:  req.Amount,
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
	c.Status(http.StatusOK)
}
