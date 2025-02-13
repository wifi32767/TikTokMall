package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/user"
)

type useridInput struct {
	Userid uint32 `json:"userid" binding:"required"`
}

type registerInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type updateInput struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type errorReturn struct {
	Error string
}

//	@Summary		注册账户
//	@Description	注册一个新的账户
//	@Tags			user
//	@Produce		json
//	@Param			input	body		registerInput	true	"注册信息"
//	@Success		200		{object}	useridInput		"用户id"
//	@Failure		400		{object}	errorReturn		"请求格式错误"
//	@Failure		500		{object}	errorReturn		"服务器错误"
//	@Router			/user/register [post]
func UserRegister(c *gin.Context) {
	input := registerInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.UserClient.Register(c.Request.Context(), &user.RegisterReq{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userid": resp.GetUserId(),
	})
}

//	@Summary		登录
//	@Description	登录一个账户
//	@Tags			user
//	@Produce		json
//	@Param			input	body	registerInput	true	"登录信息"
//	@Success		200		"成功"
//	@Failure		400		{object}	errorReturn	"请求格式错误"
//	@Failure		500		{object}	errorReturn	"服务器错误"
//	@Router			/user/login [post]
func UserLogin(c *gin.Context) {
	input := registerInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := rpc.UserClient.Login(c.Request.Context(), &user.LoginReq{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	token, err := rpc.AuthClient.DeliverToken(c.Request.Context(), &auth.DeliverTokenReq{
		UserId: resp.GetUserId(),
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.SetCookie("token", token.GetToken(), 30*24*60*60, "/", "localhost", false, true)
	c.Status(http.StatusOK)
}

//	@Summary		登出
//	@Description	登出一个账户
//	@Tags			user
//	@Produce		json
//	@Success		200	"成功"
//	@Failure		400	{object}	errorReturn	"请求格式错误"
//	@Failure		500	{object}	errorReturn	"服务器错误"
//	@Router			/user/logout [post]
func UserLogout(c *gin.Context) {
	_, err := rpc.AuthClient.DeleteToken(c.Request.Context(), &auth.DeleteTokenReq{
		Token: c.MustGet("token").(string),
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.Status(http.StatusOK)
}

//	@Summary		删除账户
//	@Description	删除一个账户
//	@Tags			user
//	@Produce		json
//	@Success		200	"成功"
//	@Failure		400	{object}	errorReturn	"请求格式错误"
//	@Failure		500	{object}	errorReturn	"服务器错误"
//	@Router			/user/delete [delete]
func UserDelete(c *gin.Context) {
	input := registerInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := rpc.UserClient.Delete(c.Request.Context(), &user.DeleteReq{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	if !ok.GetSuccess() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "删除失败"})
		return
	}
	c.Status(http.StatusOK)
}

//	@Summary		修改密码
//	@Description	修改一个账户的密码
//	@Tags			user
//	@Produce		json
//	@Param			input	body	updateInput	true	"修改信息"
//	@Success		200		"成功"
//	@Failure		400		{object}	errorReturn	"请求格式错误"
//	@Failure		500		{object}	errorReturn	"服务器错误"
//	@Router			/user/update [put]
func UserUpdate(c *gin.Context) {
	input := updateInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := rpc.UserClient.Update(c.Request.Context(), &user.UpdateReq{
		Username:    input.Username,
		OldPassword: input.OldPassword,
		NewPassword: input.NewPassword,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	if !ok.GetSuccess() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "修改失败"})
		return
	}
	_, err = rpc.AuthClient.DeleteAllTokens(c.Request.Context(), &auth.DeleteAllTokensReq{
		UserId: c.MustGet("userid").(uint32),
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.Status(http.StatusOK)
}
