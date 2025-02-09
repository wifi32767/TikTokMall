package main

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/wifi32767/TikTokMall/app/cart/biz/dal"
	"github.com/wifi32767/TikTokMall/app/cart/biz/model"
	"github.com/wifi32767/TikTokMall/app/cart/infra"
	cart "github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/product"
)

// CartServiceImpl implements the last service interface defined in the IDL.
type CartServiceImpl struct{}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.Empty, err error) {
	klog.Infof("AddItem: %v", req)
	// 首先要确认有这个物品
	_, err = infra.ProductClient.GetProduct(ctx, &product.GetProductReq{
		Id: req.GetItem().GetProductId(),
	})
	if err != nil {
		klog.Error(err)
		return
	}
	item := &model.CartItem{
		UserId:    req.GetUserId(),
		ProductId: req.GetItem().GetProductId(),
		Quantity:  req.GetItem().GetQuantity(),
	}
	err = model.AddItem(dal.DB, ctx, item)
	return
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	klog.Infof("GetCart: %v", req)
	ct, err := model.GetCart(dal.DB, ctx, req.GetUserId())
	if err != nil {
		klog.Error(err)
		return
	}
	items := make([](*cart.CartItem), 0)
	for _, item := range *ct {
		if item.Quantity == 0 {
			continue
		}
		items = append(items, &cart.CartItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		})
	}
	resp = &cart.GetCartResp{
		Cart: &cart.Cart{
			UserId: req.GetUserId(),
			Items:  items,
		},
	}
	return
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.Empty, err error) {
	klog.Infof("EmptyCart: %v", req)
	err = model.EmptyCart(dal.DB, ctx, req.GetUserId())
	return
}
