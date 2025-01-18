package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/wifi32767/TikTokMall/backend/conf"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth/authservice"
)

var AuthClient authservice.Client

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
}
