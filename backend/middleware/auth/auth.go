package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/conf"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/user"
)

// 普通路由要登录后才能访问
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		resp, err := rpc.AuthClient.VerifyToken(c.Request.Context(), &auth.VerifyTokenReq{
			Token: token,
		})
		if err != nil {
			utils.ErrorHandle(c, err)
			c.Abort()
			return
		}
		if !resp.GetRes() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		c.Set("userid", uint(resp.GetUserId()))
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

// 保护名单内的路由经过鉴权方可访问
// 要先经过token验证
func ProtectedAuthentication() gin.HandlerFunc {
	conf.GetConf()
	return func(c *gin.Context) {
		if conf.ProtectedListRe.MatchString(c.Request.URL.Path) {
			userId := c.GetUint("userid")
			if userId == 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
				c.Abort()
				return
			}
			resp, err := rpc.UserClient.GetUserPermission(c.Request.Context(), &user.GetUserPermissionReq{
				UserId: uint32(userId),
			})
			if err != nil {
				utils.ErrorHandle(c, err)
				c.Abort()
				return
			}
			fmt.Println("aaa", resp)
			if resp.GetPermission() != 2 {
				c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
