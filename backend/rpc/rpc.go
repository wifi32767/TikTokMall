package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/wifi32767/TikTokMall/backend/conf"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth/authservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart/cartservice"
	productservice "github.com/wifi32767/TikTokMall/rpc/kitex_gen/product/productcatalogservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/user/userservice"
)

var (
	AuthClient    authservice.Client
	UserClient    userservice.Client
	ProductClient productservice.Client
	CartClient    cartservice.Client
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
}
