package biz

import (
	"context"
	"fmt"
	"time"

	ak "github.com/seth16888/wxtoken/api/v1"
	"go.uber.org/zap"
)

type AccessTokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint64 `json:"expires_in"`
}

type AccessTokenUsecase struct {
	appRepo AppRepo
	log     *zap.Logger
	cli     ak.TokenClient
	timeout time.Duration
}

// func (a *AccessTokenUsecase) GetToken(ctx context.Context, appId string) (*AccessTokenRes, error) {
//   // 获取公众号信息
// 	mpInfo, err := a.appRepo.Get(ctx, appId)
// 	if err != nil {
// 		return nil, fmt.Errorf("get app error: %s", err.Error())
// 	}

//   res, err := a.FetchAccessToken(ctx, appId, mpInfo.MpId)
//   if err != nil {
//     return nil, err
//   }
//   return res, nil
// }

func NewAccessTokenUsecase(cli ak.TokenClient, logger *zap.Logger,
	timeout time.Duration, appRepo AppRepo,
) *AccessTokenUsecase {
	return &AccessTokenUsecase{cli: cli, log: logger, timeout: timeout, appRepo: appRepo}
}

func (a *AccessTokenUsecase) FetchAccessToken(ctx context.Context,
	appid string, mpId string,
) (*AccessTokenRes, error) {
	// 设置超时时间
	c, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	a.log.Debug("fetch access token", zap.String("appid", appid), zap.String("mpid", mpId))
	req := &ak.GetTokenRequest{
		AppId: appid,
		MpId:  mpId,
	}
	res, err := a.cli.GetAccessToken(c, req)
	if err != nil {
		a.log.Error("fetch access token error", zap.Error(err), zap.String("appid", appid), zap.String("mpid", mpId))
		return nil, fmt.Errorf("fetch access token error: %s", err.Error())
	}
	a.log.Debug("-> access token ok", zap.String("mpid", mpId), zap.String("access_token", res.AccessToken), zap.Uint64("expires_in", res.ExpiresIn))

	return &AccessTokenRes{
		AccessToken: res.AccessToken,
		ExpiresIn:   res.ExpiresIn,
	}, nil
}
