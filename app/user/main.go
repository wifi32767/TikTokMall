package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/wifi32767/TikTokMall/app/user/biz/dal"
	"github.com/wifi32767/TikTokMall/app/user/conf"
	user "github.com/wifi32767/TikTokMall/rpc/kitex_gen/user/userservice"
)

func main() {
	// log
	klog.SetLevel(conf.LogLevel())
	// mysql
	dal.MysqlInit()
	// kitex
	opts := kitexInit()
	svr := user.NewServer(new(UserServiceImpl), opts...)

	err := svr.Run()

	if err != nil {
		klog.Error(err.Error())
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
