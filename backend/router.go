package main

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/wifi32767/TikTokMall/backend/biz/handler"
	"github.com/wifi32767/TikTokMall/backend/biz/handler/test_handler"
	_ "github.com/wifi32767/TikTokMall/backend/docs"
	"github.com/wifi32767/TikTokMall/backend/middleware"
)

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func Register(r *gin.Engine) {

	// test下的路由为了方便测试，直接尽可能多的用get
	test := r.Group("/test")
	test.GET("/auth/deliver", test_handler.DeliverToken)
	test.GET("/auth/verify", test_handler.VerifyToken)
	test.GET("/auth/delete", test_handler.DeleteToken)
	test.GET("/auth/deleteall", test_handler.DeleteAllTokens)

	api := r.Group("/api")

	user := api.Group("/user")
	user.POST("/register", handler.UserRegister)
	user.POST("/login", handler.UserLogin)
	user.POST("/logout", middleware.Authentication(), handler.UserLogout)
	user.DELETE("/delete", middleware.Authentication(), handler.UserDelete)
	user.PUT("/update", middleware.Authentication(), handler.UserUpdate)

	r.GET("/ping", Ping)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
