package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/conf"
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

// 在白名单内的路由无需登录即可访问
func WhiteListAuthentication() gin.HandlerFunc {
	conf.GetConf()
	return func(c *gin.Context) {
		if conf.WhiteListRe.MatchString(c.Request.URL.Path) {
			c.Next()
		} else {
			Authentication()(c)
		}
	}
}
