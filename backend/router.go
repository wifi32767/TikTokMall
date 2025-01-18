package main

import (
	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/biz/handler"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Register(r *gin.Engine) {
	test := r.Group("/test")
	test.GET("/auth/deliver", handler.DeliverTokenByRPC)
	test.GET("/auth/verify", handler.VerifyTokenByRPC)
	test.GET("/auth/delete", handler.DeleteTokenByRPC)

	// api := r.Group("/api")

	r.GET("/ping", Ping)
}
