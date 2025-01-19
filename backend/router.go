package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/wifi32767/TikTokMall/backend/biz/handler"
	_ "github.com/wifi32767/TikTokMall/backend/docs"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Register(r *gin.Engine) {

	// test下的路由为了方便测试，直接尽可能多的用get
	test := r.Group("/test")
	test.GET("/auth/deliver", handler.DeliverToken)
	test.GET("/auth/verify", handler.VerifyToken)
	test.GET("/auth/delete", handler.DeleteToken)
	test.GET("/auth/deleteall", handler.DeleteAllTokens)

	// api := r.Group("/api")

	r.GET("/ping", Ping)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
