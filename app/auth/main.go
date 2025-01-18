package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/wifi32767/TikTokMall/app/auth/conf"
	"github.com/wifi32767/TikTokMall/app/auth/dal"
	log "github.com/wifi32767/TikTokMall/common/logger"
	auth "github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth/authservice"
)

func main() {
	// log
	log.Init(conf.GetConf().Kitex.Log_level)
	// redis
	dal.RedisInit()
	// kitex
	opts := kitexInit()
	svr := auth.NewServer(new(AuthServiceImpl), opts...)

	err := svr.Run()

	if err != nil {
		log.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1"+conf.GetConf().Kitex.Address)
	opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	// consul
	r, err := consul.NewConsulRegister("127.0.0.1:8500")
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
