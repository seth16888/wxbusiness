package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/seth16888/wxcommon/mp"
	ap "github.com/seth16888/wxproxy/api/v1"
	"go.uber.org/zap"
)

type APIProxyUsecase struct {
	cli     ap.MpproxyClient
	log     *zap.Logger
	timeout time.Duration
  tokenCli *AccessTokenUsecase
}

func NewAPIProxyUsecase(cli ap.MpproxyClient, log *zap.Logger,
  timeout time.Duration, tokenUc  *AccessTokenUsecase) *APIProxyUsecase {
	return &APIProxyUsecase{
		cli:     cli,
		log:     log,
    timeout: timeout,
    tokenCli: tokenUc,
	}
}

func (a *APIProxyUsecase) GetAccessToken(ctx context.Context,
	appId, mpId string) (string, error) {
	// 获取access token
	token, err := a.tokenCli.FetchAccessToken(ctx, appId, mpId)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// CreateMenu 创建菜单
func (a *APIProxyUsecase) CreateMenu(ctx context.Context,
	token string, params *mp.CreateMenuReq) error {
	// 转换请求参数
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	req := ap.CreateMenuRequest{
		AccessToken: token,
		Button:      convertButtons(params.Button),
		Matchrule:   convertMatchRule(params.MatchRule),
	}
	if req.Button == nil {
		return fmt.Errorf("button cannot be nil")
	}

  // 设置超时时间
  c, cancel := context.WithTimeout(ctx, 15 *time.Second)
  defer cancel()
  // 调用API
  a.log.Debug("call api: create menu", zap.String("access_token", token))
	rt, err := a.cli.CreateMenu(c, &req)
	if err != nil {
    a.log.Error("create menu error", zap.Error(err))
		return fmt.Errorf("api error: %s", err.Error())
	}
	if rt.Errcode != 0 {
		a.log.Error("create menu error", zap.Int64("code", rt.Errcode),
			zap.String("msg", rt.Errmsg))
    return fmt.Errorf("call api error: %s", rt.Errmsg)
	}

	return nil
}

// convertButtons converts []*mp.Button to []*ap.MenuButton
func convertButtons(buttons []*mp.Button) []*ap.MenuButton {
	if buttons == nil {
		return nil
	}
	converted := make([]*ap.MenuButton, len(buttons))
	for i, btn := range buttons {
		converted[i] = &ap.MenuButton{
			Name:      btn.Name,
			Type:      btn.Type,
			Key:       btn.Key,
			Url:       btn.URL,
			MediaId:   btn.MediaID,
			AppId:     btn.AppID,
			PagePath:  btn.PagePath,
			SubButton: convertButtons(btn.SubButtons),
		}
	}
	return converted
}

// convertMatchRule converts *mp.MatchRule to *ap.ConditionalMatchRule
func convertMatchRule(rule *mp.MatchRule) *ap.ConditionalMatchRule {
	if rule == nil {
		return nil
	}
	return &ap.ConditionalMatchRule{
		TagId:              rule.TagId,
		ClientPlatformType: rule.ClientPlatformType,
	}
}
