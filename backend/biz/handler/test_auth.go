package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth"
)

func DeliverTokenByRPC(c *gin.Context) {
	var input struct {
		Userid int32 `json:"userid"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if input.Userid == 0 {
		c.JSON(400, gin.H{"error": "userid must be set"})
		return
	}
	resp, err := rpc.AuthClient.DeliverTokenByRPC(c.Request.Context(), &auth.DeliverTokenReq{
		UserId: input.Userid,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"token": resp.Token,
	})
}

func VerifyTokenByRPC(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.AuthClient.VerifyTokenByRPC(c.Request.Context(), &auth.VerifyTokenReq{
		Token: input.Token,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"res":    resp.Res,
		"userid": resp.UserId,
	})
}

func DeleteTokenByRPC(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	token, err := utils.ParseToken(input.Token)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	_, err = rpc.AuthClient.DeleteTokenByRPC(c.Request.Context(), &auth.DeleteTokenReq{
		UserId: token.Userid,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{})
}
