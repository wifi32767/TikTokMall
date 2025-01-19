package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth"
)

func DeliverToken(c *gin.Context) {
	var input struct {
		Userid int32 `json:"userid"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Userid == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must be set"})
		return
	}
	resp, err := rpc.AuthClient.DeliverToken(c.Request.Context(), &auth.DeliverTokenReq{
		UserId: input.Userid,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": resp.Token,
	})
}

func VerifyToken(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token must be set"})
		return
	}
	resp, err := rpc.AuthClient.VerifyToken(c.Request.Context(), &auth.VerifyTokenReq{
		Token: input.Token,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"res":    resp.Res,
		"userid": resp.UserId,
	})
}

func DeleteToken(c *gin.Context) {
	var input struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token must be set"})
		return
	}
	_, err := rpc.AuthClient.DeleteToken(c.Request.Context(), &auth.DeleteTokenReq{
		Token: input.Token,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func DeleteAllTokens(c *gin.Context) {
	var input struct {
		Userid int32 `json:"userid"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Userid == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userid must be set"})
		return
	}
	_, err := rpc.AuthClient.DeleteAllTokens(c.Request.Context(), &auth.DeleteAllTokensReq{
		UserId: input.Userid,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
