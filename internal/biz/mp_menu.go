package biz

import (
	"context"

	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxcommon/mp"
)

type MPMenuUsecase struct {
	repo     PlatformAppRepo
	tokenUc  *AccessTokenUsecase
	apiProxy *APIProxyUsecase
}

func NewMPMenuUsecase(repo PlatformAppRepo,
	tokenUc *AccessTokenUsecase,
	apiProxy *APIProxyUsecase,
) *MPMenuUsecase {
	return &MPMenuUsecase{repo: repo, tokenUc: tokenUc, apiProxy: apiProxy}
}

func (u *MPMenuUsecase) Pull(ctx context.Context, appId string) *r.R {
	// 获取公众号信息
	mpInfo, err := u.repo.Get(ctx, appId)
	if err != nil {
		return r.Error(400, "app not found")
	}
	// accessToken
	_, err = u.tokenUc.FetchAccessToken(ctx, appId, mpInfo.AppId)
	if err != nil {
		return r.Error(403, "fetch access token error")
	}

	return r.SuccessData(mpInfo)
}

func (u *MPMenuUsecase) Create(ctx context.Context, pId string, params *mp.CreateMenuReq) *r.R {
	// 获取公众号信息
	mpInfo, err := u.repo.Get(ctx, pId)
	if err != nil {
		return r.Error(400, "app not found")
	}
	// accessToken
	akRes, err := u.tokenUc.FetchAccessToken(ctx, pId, mpInfo.AppId)
	if err != nil {
		return r.Error(403, "fetch access token error")
	}

	err = u.apiProxy.CreateMenu(ctx, akRes.AccessToken, params)
	if err != nil {
		return r.Error(400, "create menu error")
	}

	return r.Success()
}
