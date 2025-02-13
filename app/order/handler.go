package main

import (
	"context"
	"net/http"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"github.com/wifi32767/TikTokMall/app/order/biz/dal"
	"github.com/wifi32767/TikTokMall/app/order/biz/model"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart"
	order "github.com/wifi32767/TikTokMall/rpc/kitex_gen/order"
)

// OrderServiceImpl implements the last service interface defined in the IDL.
type OrderServiceImpl struct{}

// PlaceOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	klog.Infof("PlaceOrder: %v", req)
	orderid, _ := uuid.NewUUID()
	o := &model.Order{
		OrderId:      orderid.String(),
		UserId:       req.GetUserId(),
		UserCurrency: req.GetUserCurrency(),
		Consignee: model.Consignee{
			Email: req.GetEmail(),
		},
		OrderState: model.OrderStatePlaced,
	}
	if req.Address != nil {
		a := req.GetAddress()
		o.Consignee.Country = a.Country
		o.Consignee.State = a.State
		o.Consignee.City = a.City
		o.Consignee.StreetAddress = a.StreetAddress
	}
	orderItemList := make([]model.OrderItem, 0)
	for _, item := range req.GetOrderItems() {
		orderItemList = append(orderItemList, model.OrderItem{
			OrderIdRefer: o.OrderId,
			ProductId:    item.GetItem().GetProductId(),
			Quantity:     item.GetItem().GetQuantity(),
			Cost:         item.GetCost(),
		})
	}
	err = model.PlaceOrder(dal.DB, ctx, o, &orderItemList)
	if err != nil {
		klog.Error(err)
		return
	}
	resp = &order.PlaceOrderResp{
		Order: &order.OrderResult{
			OrderId: o.OrderId,
		},
	}
	return
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	klog.Infof("ListOrder: %v", req)
	orders, err := model.ListOrder(dal.DB, ctx, req.GetUserId())
	if err != nil {
		klog.Error(err)
		return
	}
	resp = &order.ListOrderResp{
		Orders: make([]*order.Order, 0),
	}
	for _, o := range *orders {
		items := make([]*order.OrderItem, 0)
		for _, item := range o.OrderItems {
			items = append(items, &order.OrderItem{
				Item: &cart.CartItem{
					ProductId: item.ProductId,
					Quantity:  item.Quantity,
				},
				Cost: item.Cost,
			})
		}
		resp.Orders = append(resp.Orders, &order.Order{
			OrderId:      o.OrderId,
			UserId:       o.UserId,
			UserCurrency: o.UserCurrency,
			Address: &order.Address{
				StreetAddress: o.Consignee.StreetAddress,
				City:          o.Consignee.City,
				State:         o.Consignee.State,
				Country:       o.Consignee.Country,
				ZipCode:       o.Consignee.ZipCode,
			},
			Email:      o.Consignee.Email,
			OrderItems: items,
			CreatedAt:  int32(o.CreatedAt.Second()),
			State:      string(o.OrderState),
		})
	}
	return
}

// UpdateOrderState implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) UpdateOrderState(ctx context.Context, req *order.UpdateOrderStateReq) (resp *order.Empty, err error) {
	klog.Infof("UpdateOrderState: %v", req)
	var state model.OrderState
	switch req.GetState() {
	case 0:
		state = model.OrderStatePlaced
	case 1:
		state = model.OrderStatePaid
	case 2:
		state = model.OrderStateCanceled
	default:
		err = kerrors.NewBizStatusError(http.StatusBadRequest, "invalid state")
		return
	}
	err = model.UpdateOrderState(dal.DB, ctx, req.GetUserId(), req.GetOrderId(), state)
	return
}

// GetOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) GetOrder(ctx context.Context, req *order.GetOrderReq) (resp *order.GetOrderResp, err error) {
	klog.Infof("GetOrder: %v", req)
	o, err := model.GetOrder(dal.DB, ctx, req.GetOrderId())
	if err != nil {
		klog.Error(err)
		return
	}
	items := make([]*order.OrderItem, 0)
	for _, item := range o.OrderItems {
		items = append(items, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			},
			Cost: item.Cost,
		})
	}
	resp = &order.GetOrderResp{
		Order: &order.Order{
			OrderId:      o.OrderId,
			UserId:       o.UserId,
			UserCurrency: o.UserCurrency,
			Address: &order.Address{
				StreetAddress: o.Consignee.StreetAddress,
				City:          o.Consignee.City,
				State:         o.Consignee.State,
				Country:       o.Consignee.Country,
				ZipCode:       o.Consignee.ZipCode,
			},
			Email:      o.Consignee.Email,
			OrderItems: items,
			CreatedAt:  int32(o.CreatedAt.Second()),
			State:      string(o.OrderState),
		},
	}
	return
}
