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
	log     *zap.Logger
	cli     ak.TokenClient
	timeout time.Duration
}

func NewAccessTokenUsecase(cli ak.TokenClient, logger *zap.Logger,
  timeout time.Duration) *AccessTokenUsecase {
	return &AccessTokenUsecase{cli: cli, log: logger, timeout: timeout}
}

func (a *AccessTokenUsecase) FetchAccessToken(ctx context.Context, appid string, mpId string) (*AccessTokenRes, error) {
  // 设置超时时间
  c, cancel := context.WithTimeout(ctx, a.timeout)
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
