package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/wifi32767/TikTokMall/app/checkout/conf"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart/cartservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order/orderservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/payment/paymentservice"
	productservice "github.com/wifi32767/TikTokMall/rpc/kitex_gen/product/productcatalogservice"
)

var (
	ProductClient productservice.Client
	CartClient    cartservice.Client
	OrderClient   orderservice.Client
	PaymentClient paymentservice.Client
)

func Init() {
	r, err := consul.NewConsulResolver(conf.GetConf().Kitex.ConsulAddress)
	if err != nil {
		panic(err)
	}

	ProductClient = productservice.MustNewClient(
		"product",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
	CartClient = cartservice.MustNewClient(
		"cart",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
	OrderClient = orderservice.MustNewClient(
		"order",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
	PaymentClient = paymentservice.MustNewClient(
		"payment",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
	klog.Infof("rpc init success")
}
