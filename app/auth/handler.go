package main

import (
	"context"
	"strconv"
	"time"

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
	// 生成一个token，存入redis
	token, err := utils.GenerateToken(req.GetUserId())
	if err != nil {
		return nil, err
	}
	err = dal.RedisDB.Set(ctx, strconv.Itoa(int(req.GetUserId())), token, oneYear).Err()
	if err != nil {
		return
	}

	resp = &auth.DeliveryResp{
		Token: token,
	}
	return
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	// 验证redis中是否有这个token
	token, err := utils.ParseToken(req.Token)
	if err != nil {
		return nil, err
	}
	t, err := dal.RedisDB.Get(ctx, strconv.Itoa(int(token.Userid))).Result()
	resp = &auth.VerifyResp{
		Res: true,
	}
	if err != nil {
		if err == redis.Nil {
			err = nil
			resp.Res = false
		}
	} else {
		if t != req.Token {
			resp.Res = false
		}
		resp.UserId = token.Userid
	}
	return
}

// DeleteTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeleteTokenByRPC(ctx context.Context, req *auth.DeleteTokenReq) (resp *auth.Empty, err error) {
	err = dal.RedisDB.Del(ctx, strconv.Itoa(int(req.GetUserId()))).Err()
	return
}
