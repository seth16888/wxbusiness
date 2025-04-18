package biz

import (
	"fmt"
	"net/url"

	"github.com/seth16888/wxcommon/domain"
	"github.com/seth16888/wxcommon/paths"
	v1 "github.com/seth16888/wxproxy/api/v1"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/internal/model/response"
)

type MpQRCodeUsecase struct {
	log      *zap.Logger
	repo     AppRepo
	tokenUc  *AccessTokenUsecase
	apiProxy *APIProxyUsecase
}

// GetURL 获取二维码URL
func (m *MpQRCodeUsecase) GetURL(c *gin.Context, ticket string) string {
	encodedTicket := url.QueryEscape(ticket)

	url := fmt.Sprintf("https://%s%s?ticket=%s",
		domain.GetMPDomain(),
		paths.Path_Get_QRCode_URL,
		encodedTicket,
	)

	return url
}

// CreateLimit 永久二维码，是无过期时间的，但数量较少（目前为最多10万个）。
// 永久二维码主要用于适用于账号绑定、用户来源统计等场景。
func (m *MpQRCodeUsecase) CreateLimit(c *gin.Context, appId string,
	req *request.CreateQRCodeReq,
) (*response.Ticket, error) {
  app, err := GetAppInfoFromCtx(c)
	if err != nil {
		m.log.Error("get app info error", zap.Error(err))
		return nil, fmt.Errorf("get app info error")
	}
	mpId := app.MpId
  // 获取access_token
  token,err:= m.tokenUc.FetchAccessToken(c, appId, mpId)
  if err!=nil{
    m.log.Error("GetToken error", zap.Error(err))
    return nil,err
  }

	params := v1.CreateQRCodeRequest{
		AccessToken:   token.AccessToken,
		ExpireSeconds: 0,
		Scene:         req.Scene,
	}
	res, err := m.apiProxy.cli.CreateLimitQRCode(c, &params)
	if err != nil {
		m.log.Error("CreateLimit error", zap.Error(err))
		return nil, err
	}

	resp := response.Ticket{
		Ticket:        res.Ticket,
		URL:           res.URL,
		ExpireSeconds: res.ExpireSeconds,
	}

	return &resp, nil
}

// CreateTemporary 临时二维码，
// 是有过期时间的，最长可以设置为在二维码生成后的30天（即2592000秒）后过期，
// 但能够生成较多数量。临时二维码主要用于账号绑定等不要求二维码永久保存的业务场景。
func (m *MpQRCodeUsecase) CreateTemporary(c *gin.Context, appId string,
	req *request.CreateQRCodeReq,
) (*response.Ticket, error) {
  app, err := GetAppInfoFromCtx(c)
	if err != nil {
		m.log.Error("get app info error", zap.Error(err))
		return nil, fmt.Errorf("get app info error")
	}
	mpId := app.MpId
  // 获取access_token
  token,err:= m.tokenUc.FetchAccessToken(c, appId, mpId)
  if err!=nil{
    m.log.Error("GetToken error", zap.Error(err))
    return nil,err
  }

	params := v1.CreateQRCodeRequest{
		AccessToken:   token.AccessToken,
		ExpireSeconds: int64(req.Exp),
		Scene:         req.Scene,
	}
	res, err := m.apiProxy.cli.CreateTemporaryQRCode(c, &params)
	if err != nil {
		m.log.Error("CreateTemporaryQRCode error", zap.Error(err))
		return nil, err
	}

	resp := response.Ticket{
		Ticket:        res.Ticket,
		URL:           res.URL,
		ExpireSeconds: res.ExpireSeconds,
	}

	return &resp, nil
}

func NewMpQRCodeUsecase(
	log *zap.Logger,
	repo AppRepo,
	tokenUc *AccessTokenUsecase,
	apiProxy *APIProxyUsecase,
) *MpQRCodeUsecase {
	return &MpQRCodeUsecase{
		log:      log,
		repo:     repo,
		tokenUc:  tokenUc,
		apiProxy: apiProxy,
	}
}
