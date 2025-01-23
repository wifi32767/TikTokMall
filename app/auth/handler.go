package main

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"
	"github.com/wifi32767/TikTokMall/app/auth/dal"
	"github.com/wifi32767/TikTokMall/app/auth/utils"
	auth "github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

var oneMonth = 30 * 24 * 60 * 60 * time.Second

// DeliverToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverToken(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	klog.Infof("DeliverToken: %v", req)
	resp = &auth.DeliveryResp{}
	// 生成一个token，存入redis
	token, err := utils.GenerateToken(req.GetUserId())
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	// redis 7.4.0之后支持为Hash中的各个字段分别设置过期时间
	// 因此需要用比较新的redis库，9.7.0之上
	err = dal.RedisClient.HSet(ctx, strconv.Itoa(int(req.GetUserId())), token, 1).Err()
	if err != nil {
		klog.Error(err)
		return
	}
	err = dal.RedisClient.HExpire(ctx, strconv.Itoa(int(req.GetUserId())), oneMonth, token).Err()
	if err != nil {
		klog.Error(err)
		return
	}
	resp.Token = token
	return
}

// VerifyToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyToken(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	klog.Infof("VerifyToken: %v", req)
	resp = &auth.VerifyResp{}
	// 验证redis中是否有这个token
	token, err := utils.ParseToken(req.Token)
	// token解析失败，直接确认token无效
	if err != nil {
		resp.Res = false
		err = nil
		return
	}
	err = dal.RedisClient.HGet(ctx, strconv.Itoa(int(token.Userid)), req.GetToken()).Err()
	if err != nil {
		// 没有这个token
		if err == redis.Nil {
			err = nil
			resp.Res = false
		} else {
			klog.Error(err)
		}
	} else {
		resp.Res = true
		resp.UserId = token.Userid
		// 续期token
		dal.RedisClient.HExpire(ctx, strconv.Itoa(int(token.Userid)), oneMonth, req.GetToken())
	}
	return
}

// DeleteToken implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeleteToken(ctx context.Context, req *auth.DeleteTokenReq) (resp *auth.Empty, err error) {
	klog.Infof("DeleteToken: %v", req)
	token, err := utils.ParseToken(req.GetToken())
	if err != nil {
		err = nil
		return
	}
	// 删除redis中的token
	err = dal.RedisClient.HDel(ctx, strconv.Itoa(int(token.Userid)), req.GetToken()).Err()
	if err == redis.Nil {
		err = nil
	} else {
		klog.Error(err)
	}
	return
}

// DeleteAllTokens implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeleteAllTokens(ctx context.Context, req *auth.DeleteAllTokensReq) (resp *auth.Empty, err error) {
	klog.Infof("DeleteAllTokens: %v", req)
	// 删除redis中的token
	err = dal.RedisClient.Del(ctx, strconv.Itoa(int(req.GetUserId()))).Err()
	if err == redis.Nil {
		err = nil
	} else {
		klog.Error(err)
	}
	return
}
