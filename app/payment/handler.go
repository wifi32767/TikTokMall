package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	creditcard "github.com/durango/go-credit-card"
	"github.com/google/uuid"
	"github.com/wifi32767/TikTokMall/app/payment/biz/dal"
	"github.com/wifi32767/TikTokMall/app/payment/biz/model"
	payment "github.com/wifi32767/TikTokMall/rpc/kitex_gen/payment"
)

// PaymentServiceImpl implements the last service interface defined in the IDL.
type PaymentServiceImpl struct{}

// Charge implements the PaymentServiceImpl interface.
func (s *PaymentServiceImpl) Charge(ctx context.Context, req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	klog.Infof("Charge: %v", req)
	card := creditcard.Card{
		Number: req.CreditCard.CreditCardNumber,
		Cvv:    strconv.Itoa(int(req.CreditCard.CreditCardCvv)),
		Month:  strconv.Itoa(int(req.CreditCard.CreditCardExpirationMonth)),
		Year:   strconv.Itoa(int(req.CreditCard.CreditCardExpirationYear)),
	}
	// 检测信用卡信息是否合法
	err = card.Validate(true)
	if err != nil {
		err = kerrors.NewBizStatusError(http.StatusBadRequest, err.Error())
		return
	}
	// 生成交易ID
	translationId, err := uuid.NewRandom()
	if err != nil {
		err = kerrors.NewBizStatusError(http.StatusInternalServerError, err.Error())
		klog.Error(err)
		return
	}
	// 创建支付记录
	err = model.CreatePaymentLog(dal.DB, ctx, &model.PaymentLog{
		UserId:        req.UserId,
		OrderId:       req.OrderId,
		TransactionId: translationId.String(),
		Amount:        req.Amount,
		PayAt:         time.Now(),
	})
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	resp = &payment.ChargeResp{TransactionId: translationId.String()}
	return
}
