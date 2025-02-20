package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/wifi32767/TikTokMall/app/checkout/infra/rpc"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart"
	checkout "github.com/wifi32767/TikTokMall/rpc/kitex_gen/checkout"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/payment"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/product"
)

// CheckoutServiceImpl implements the last service interface defined in the IDL.
type CheckoutServiceImpl struct{}

// Checkout implements the CheckoutServiceImpl interface.
// 1. 获取购物车信息
// 2. 计算购物车总额
// 3. 创建订单
// 4. 清空购物车
// 5. 支付
// 6. 更改订单状态
func (s *CheckoutServiceImpl) Checkout(ctx context.Context, req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	klog.Infof("Checkout request: %v", req)
	// 1. 获取购物车信息
	cartRes, err := rpc.CartClient.GetCart(ctx, &cart.GetCartReq{
		UserId: req.GetUserId(),
	})
	if err != nil {
		return
	}
	// 2. 计算购物车总额
	orderList := make([]*order.OrderItem, 0)
	totalCost := float32(0)
	for _, item := range cartRes.GetCart().GetItems() {
		prod, prodErr := rpc.ProductClient.GetProduct(ctx, &product.GetProductReq{
			Id: item.GetProductId(),
		})
		if prodErr != nil {
			err = prodErr
			return
		}
		cost := float32(prod.GetProduct().GetPrice()) * float32(item.GetQuantity())
		totalCost += cost
		orderList = append(orderList, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId: item.GetProductId(),
				Quantity:  item.GetQuantity(),
			},
			Cost: cost,
		})
	}
	// 3. 创建订单
	orderReq := &order.PlaceOrderReq{
		UserId:       req.GetUserId(),
		UserCurrency: "RMB",
		OrderItems:   orderList,
		Email:        req.GetEmail(),
	}
	if req.Address != nil {
		orderReq.Address = req.GetAddress()
	}
	orderRes, err := rpc.OrderClient.PlaceOrder(ctx, orderReq)
	if err != nil {
		return
	}
	// 4. 清空购物车
	_, err = rpc.CartClient.EmptyCart(ctx, &cart.EmptyCartReq{
		UserId: req.GetUserId(),
	})
	if err != nil {
		return
	}
	// 5. 支付
	payRes, err := rpc.PaymentClient.Charge(ctx, &payment.ChargeReq{
		UserId:     req.GetUserId(),
		OrderId:    orderRes.GetOrder().GetOrderId(),
		Amount:     totalCost,
		CreditCard: req.GetCreditCard(),
	})
	if err != nil {
		return
	}
	// 6. 更改订单状态
	_, err = rpc.OrderClient.UpdateOrderState(ctx, &order.UpdateOrderStateReq{
		UserId:  req.GetUserId(),
		OrderId: orderRes.GetOrder().GetOrderId(),
		State:   1,
	})
	if err != nil {
		return
	}
	resp = &checkout.CheckoutResp{
		OrderId:       orderRes.GetOrder().GetOrderId(),
		TransactionId: payRes.GetTransactionId(),
	}
	return
}
