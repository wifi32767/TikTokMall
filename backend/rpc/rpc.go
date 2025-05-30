package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/wifi32767/TikTokMall/backend/conf"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth/authservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart/cartservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/checkout/checkoutservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order/orderservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/payment/paymentservice"
	productservice "github.com/wifi32767/TikTokMall/rpc/kitex_gen/product/productcatalogservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/user/userservice"
)

var (
	AuthClient     authservice.Client
	UserClient     userservice.Client
	ProductClient  productservice.Client
	CartClient     cartservice.Client
	OrderClient    orderservice.Client
	PaymentClient  paymentservice.Client
	CheckoutClient checkoutservice.Client
)

func Init() {
	r, err := consul.NewConsulResolver(conf.GetConf().Rpc.Consul_address)
	if err != nil {
		panic(err)
	}
	AuthClient = authservice.MustNewClient(
		"auth",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
	UserClient = userservice.MustNewClient(
		"user",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
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
	CheckoutClient = checkoutservice.MustNewClient(
		"checkout",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
}
