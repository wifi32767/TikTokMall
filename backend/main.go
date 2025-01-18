package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/conf"
	"github.com/wifi32767/TikTokMall/backend/rpc"
)

func main() {
	rpc.Init()
	r := gin.Default()
	Register(r)
	err := r.Run(conf.GetConf().Gin.Port)
	if err != nil {
		panic(err)
	}
}
