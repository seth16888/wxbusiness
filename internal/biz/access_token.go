package biz

import (
	"context"
	"fmt"

	v1 "github.com/seth16888/wxtoken/api/v1"
	"go.uber.org/zap"
)

type AccessTokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   uint64 `json:"expires_in"`
}

type AccessTokenUsecase struct {
	log *zap.Logger
	cli v1.TokenClient
}

func NewAccessTokenUsecase(cli v1.TokenClient, logger *zap.Logger) *AccessTokenUsecase {
	return &AccessTokenUsecase{cli: cli, log: logger}
}

func (a *AccessTokenUsecase) FetchAccessToken(ctx context.Context, appid string, mpId string) (*AccessTokenRes, error) {
	a.log.Debug("fetch access token", zap.String("appid", appid), zap.String("mpid", mpId))
	req := &v1.GetTokenRequest{
		AppId: appid,
		MpId:  mpId,
	}
	res, err := a.cli.GetAccessToken(ctx, req)
	if err != nil {
		a.log.Error("fetch access token error", zap.Error(err), zap.String("appid", appid), zap.String("mpid", mpId))
		return nil, fmt.Errorf("fetch access token error: %s", err.Error())
	}

	return &AccessTokenRes{
		AccessToken: res.AccessToken,
		ExpiresIn:   res.ExpiresIn,
	}, nil
}
