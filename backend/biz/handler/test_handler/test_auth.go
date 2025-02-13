package test_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth"
)

type errorReturn struct {
	Error string
}

type useridInput struct {
	Userid uint32 `json:"userid" binding:"required"`
}

type tokenInput struct {
	Token string `json:"token" binding:"required"`
}

type verifyOutput struct {
	Res    bool
	Userid uint32
}

//	@Summary		获取token
//	@Description	对于给定的userid，为它分发一个token
//	@Tags			test/auth
//	@Produce		json
//	@Param			input	body		useridInput	true	"用户id"
//	@Success		200		{object}	tokenInput	"token"
//	@Failure		400		{object}	errorReturn	"请求格式错误"
//	@Failure		500		{object}	errorReturn	"服务器错误"
//	@Router			/test/auth/deliver [get]
func DeliverToken(c *gin.Context) {
	input := useridInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

//	@Summary		验证token
//	@Description	验证一个token是否有效，以及属于哪个用户。会自动为token续期
//	@Tags			test/auth
//	@Produce		json
//	@Param			input	body		tokenInput		true	"token"
//	@Success		200		{object}	verifyOutput	"结果"
//	@Failure		400		{object}	errorReturn		"请求格式错误"
//	@Failure		500		{object}	errorReturn		"服务器错误"
//	@Router			/test/auth/verify [get]
func VerifyToken(c *gin.Context) {
	input := tokenInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

//	@Summary		删除token
//	@Description	删除一个token，即退出该用户在该机器上的登录
//	@Tags			test/auth
//	@Produce		json
//	@Param			input	body	tokenInput	true	"token"
//	@Success		200		"成功"
//	@Failure		400		{object}	errorReturn	"请求格式错误"
//	@Failure		500		{object}	errorReturn	"服务器错误"
//	@Router			/test/auth/delete [get]
func DeleteToken(c *gin.Context) {
	input := tokenInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.AuthClient.DeleteToken(c.Request.Context(), &auth.DeleteTokenReq{
		Token: input.Token,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.Status(http.StatusOK)
}

//	@Summary		删除所有token
//	@Description	删除一个用户的所有token，即退出该用户在所有机器上的登录
//	@Tags			test/auth
//	@Produce		json
//	@Param			input	body	useridInput	true	"token"
//	@Success		200		"成功"
//	@Failure		400		{object}	errorReturn	"请求格式错误"
//	@Failure		500		{object}	errorReturn	"服务器错误"
//	@Router			/test/auth/deleteall [get]
func DeleteAllTokens(c *gin.Context) {
	input := useridInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rpc.AuthClient.DeleteAllTokens(c.Request.Context(), &auth.DeleteAllTokensReq{
		UserId: input.Userid,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.Status(http.StatusOK)
}
