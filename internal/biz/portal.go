package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/seth16888/wxbusiness/internal/data/entities"
	"github.com/seth16888/wxbusiness/pkg/logger"
	wx "github.com/seth16888/wxcommon/helpers"
	"go.uber.org/zap"
)

type PortalUsecase struct {
	repo       AppRepo
	log        *zap.Logger
	tokenProxy *AccessTokenUsecase
}

func NewPortalUsecase(repo AppRepo,
  tokenProxy *AccessTokenUsecase, logger *zap.Logger) *PortalUsecase {
	return &PortalUsecase{repo: repo, log: logger, tokenProxy: tokenProxy}
}

func (p *PortalUsecase) Verify(ctx context.Context, appId string, timestamp string,
	nonce string, signature string, echostr string,
) (string, error) {
	// 获取app信息
	appInfo, err := p.repo.Get(ctx, appId)
	if err != nil {
		logger.Errorf("fetch app info error,%s: %s", err.Error(), appId)
		return "", err
	}
	// token
	token := appInfo.Token

	// 验证请求参数
	if !wx.ValidatePortalReq(timestamp, nonce, signature, token) {
		logger.Errorf("wx verify req validate failed: %s, %s, %s, %s", appId, timestamp, nonce, signature)
		return "", errors.New("wx verify req validate failed")
	}

	// 验证成功
	// 更新状态
	if appInfo.Status >= 0 { // 0-未验证,1-已暂停,2-验证成功,3-接入成功
		appInfo.Status = 2
		if err := p.repo.UpdateStatus(ctx, appInfo.ID.Hex(), appInfo.Status); err != nil {
			logger.Errorf("update app info error,%s: %s", err.Error(), appId)
			return "", err
		}

    ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    // 获取一次Access Token
		_, err := p.tokenProxy.FetchAccessToken(ctxTimeout, appId, appInfo.MpId)
		if err != nil {
			logger.Errorf("fetch access token error,%s: %s", err.Error(), appId)
			return "", err
		}
	}

	logger.Debugf("verify success: %s", appId)
	return echostr, nil
}

// GetMpApp
func (p *PortalUsecase) GetMpApp(ctx context.Context, appId string) (*entities.PlatformApp, error) {
  app, err:= p.repo.Get(ctx, appId)
  if err != nil {
    return nil, fmt.Errorf("data not found")
  }
  return app, nil
}
