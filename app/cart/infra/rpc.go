package infra

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/wifi32767/TikTokMall/app/cart/conf"
	productservice "github.com/wifi32767/TikTokMall/rpc/kitex_gen/product/productcatalogservice"
)

var ProductClient productservice.Client

func RpcInit() {
	r, err := consul.NewConsulResolver(conf.GetConf().Kitex.ConsulAddress)
	if err != nil {
		panic(err)
	}
	ProductClient = productservice.MustNewClient(
		"product",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
}
