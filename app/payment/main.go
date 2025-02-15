package main

import (
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/wifi32767/TikTokMall/app/payment/biz/dal"
	"github.com/wifi32767/TikTokMall/app/payment/conf"
	payment "github.com/wifi32767/TikTokMall/rpc/kitex_gen/payment/paymentservice"
)

func main() {

	// log
	klog.SetLevel(conf.LogLevel())
	// dal
	dal.MysqlInit()
	// kitex
	opts := kitexInit()
	svr := payment.NewServer(new(PaymentServiceImpl), opts...)

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1"+conf.GetConf().Kitex.Address)
	opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	// consul
	r, err := consul.NewConsulRegister(conf.GetConf().Kitex.Consul_address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithRegistry(r))
	opts = append(opts, server.WithRegistryInfo(&registry.Info{
		ServiceName: conf.GetConf().Kitex.Service,
		Weight:      1,
	}))
	return
}
