package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(401, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		resp, err := rpc.AuthClient.VerifyToken(c.Request.Context(), &auth.VerifyTokenReq{
			Token: token,
		})
		if err != nil {
			utils.ErrorHandle(c, err)
			return
		}
		if !resp.GetRes() {
			c.JSON(401, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		c.Set("userid", resp.GetUserId())
		c.Set("token", token)
		c.Next()
	}
}
