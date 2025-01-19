package main

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
	"github.com/wifi32767/TikTokMall/app/auth/dal"
	"github.com/wifi32767/TikTokMall/app/auth/utils"
	auth "github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

var oneYear = 365 * 24 * 60 * 60 * time.Second

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	klog.Infof("DeliverToken: %v", req)
	resp = &auth.DeliveryResp{}
	// 生成一个token，存入redis
	token, err := utils.GenerateToken(req.GetUserId())
	if err != nil {
		return nil, err
	}
	err = dal.RedisClient.Set(ctx, strconv.Itoa(int(req.GetUserId())), token, oneYear).Err()
	if err != nil {
		return
	}
	resp.Token = token
	return
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	klog.Infof("VerifyToken: %v", req)
	resp = &auth.VerifyResp{}
	// 验证redis中是否有这个token
	token, err := utils.ParseToken(req.Token)
	if err != nil {
		resp.Res = false
		return
	}
	t, err := dal.RedisClient.Get(ctx, strconv.Itoa(int(token.Userid))).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			resp.Res = false
		}
	} else {
		if t == req.Token {
			resp.Res = true
			resp.UserId = token.Userid
		}
	}
	return
}

// DeleteTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeleteTokenByRPC(ctx context.Context, req *auth.DeleteTokenReq) (resp *auth.Empty, err error) {
	klog.Infof("DeleteToken: %v", req)
	// 删除redis中的token
	err = dal.RedisClient.Del(ctx, strconv.Itoa(int(req.GetUserId()))).Err()
	return
}
