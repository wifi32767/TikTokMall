package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/wifi32767/TikTokMall/app/user/biz/dal"
	"github.com/wifi32767/TikTokMall/app/user/biz/model"
	"github.com/wifi32767/TikTokMall/app/user/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/user"
	"gorm.io/gorm"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	klog.Infof("Register: %v", req)
	// 检查用户是否已存在
	_, err = model.GetUserByUsername(dal.DB, ctx, req.GetUsername())
	if err == nil {
		err = kerrors.NewBizStatusError(http.StatusBadRequest, "user already exists")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		klog.Error(err)
		return
	}
	// 加密密码，创建用户
	pwdHash, err := utils.GetHash(req.GetPassword())
	if err != nil {
		klog.Error(err)
		return
	}
	usr, err := model.CreateUser(dal.DB, ctx, req.GetUsername(), pwdHash)
	if err != nil {
		klog.Error(err)
		return
	}
	resp = &user.RegisterResp{
		UserId: uint32(usr.Model.ID),
	}
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	klog.Infof("Login: %v", req)
	// 检查用户是否存在
	usr, err := model.GetUserByUsername(dal.DB, ctx, req.GetUsername())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = kerrors.NewBizStatusError(http.StatusNotFound, "user not found")
		return
	} else if err != nil {
		klog.Error(err)
		return
	}
	// 检查密码是否正确
	match, err := utils.CompareHash(req.GetPassword(), usr.PasswordHash)
	if err != nil {
		klog.Error(err)
		return
	}
	if !match {
		err = kerrors.NewBizStatusError(http.StatusUnauthorized, "password incorrect")
	}
	resp = &user.LoginResp{
		UserId: uint32(usr.Model.ID),
	}
	return
}

// Delete implements the UserServiceImpl interface.
func (s *UserServiceImpl) Delete(ctx context.Context, req *user.DeleteReq) (resp *user.DeleteResp, err error) {
	klog.Infof("Delete: %v", req)
	// 检查用户是否存在
	usr, err := model.GetUserByUsername(dal.DB, ctx, req.GetUsername())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = kerrors.NewBizStatusError(http.StatusNotFound, "user not found")
		return
	} else if err != nil {
		klog.Error(err)
		return
	}
	// 检查密码是否正确
	match, err := utils.CompareHash(req.GetPassword(), usr.PasswordHash)
	if err != nil {
		klog.Error(err)
		return
	}
	if !match {
		err = kerrors.NewBizStatusError(http.StatusUnauthorized, "password incorrect")
		return
	}
	// 删除用户
	err = model.DeleteUser(dal.DB, ctx, usr.Model.ID)
	resp = &user.DeleteResp{
		Success: true,
	}
	return
}

// Update implements the UserServiceImpl interface.
func (s *UserServiceImpl) Update(ctx context.Context, req *user.UpdateReq) (resp *user.UpdateResp, err error) {
	klog.Infof("Update: %v", req)
	// 检查用户是否存在
	usr, err := model.GetUserByUsername(dal.DB, ctx, req.GetUsername())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = kerrors.NewBizStatusError(http.StatusNotFound, "user not found")
		return
	} else if err != nil {
		klog.Error(err)
		return
	}
	match, err := utils.CompareHash(req.GetOldPassword(), usr.PasswordHash)
	if err != nil {
		klog.Error(err)
		return
	}
	if !match {
		err = kerrors.NewBizStatusError(http.StatusUnauthorized, "password incorrect")
		return
	}
	pwdHash, err := utils.GetHash(req.GetNewPassword())
	if err != nil {
		klog.Error(err)
		return
	}
	err = model.UpdateUser(dal.DB, ctx, usr.Model.ID, pwdHash)
	if err != nil {
		klog.Error(err)
	}
	resp = &user.UpdateResp{
		Success: true,
	}
	return
}
