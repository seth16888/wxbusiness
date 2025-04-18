package biz

import (
	"context"

	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxcommon/mp"
	"go.uber.org/zap"
)

type MPMenuUsecase struct {
	repo     AppRepo
	tokenUc  *AccessTokenUsecase
	apiProxy *APIProxyUsecase
  log       *zap.Logger
}

func NewMPMenuUsecase(repo AppRepo,
	tokenUc *AccessTokenUsecase,
	apiProxy *APIProxyUsecase,
  log       *zap.Logger,
) *MPMenuUsecase {
	return &MPMenuUsecase{repo: repo, tokenUc: tokenUc, apiProxy: apiProxy, log: log}
}

func (u *MPMenuUsecase) Pull(ctx context.Context, appId string) *r.R {
  app, err := GetAppInfoFromCtx(ctx)
	if err != nil {
		u.log.Error("get app info error", zap.Error(err))
		return r.Error(400, "get app info error")
	}
	mpId := app.MpId
	// accessToken
	_, err = u.tokenUc.FetchAccessToken(ctx, appId, mpId)
	if err != nil {
		return r.Error(403, "fetch access token error")
	}

	return r.SuccessData(nil)
}

func (u *MPMenuUsecase) Create(ctx context.Context, pId string, params *mp.CreateMenuReq) *r.R {
  app, err := GetAppInfoFromCtx(ctx)
	if err != nil {
		u.log.Error("get app info error", zap.Error(err))
		return r.Error(400, "get app info error")
	}
	mpId := app.MpId
	// accessToken
	akRes, err := u.tokenUc.FetchAccessToken(ctx, pId, mpId)
	if err != nil {
		return r.Error(403, "fetch access token error")
	}

	err = u.apiProxy.CreateMenu(ctx, akRes.AccessToken, params)
	if err != nil {
		return r.Error(400, "create menu error")
	}

	return r.Success()
}
