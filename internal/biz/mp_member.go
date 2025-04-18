package biz

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/data/entities"
	v1 "github.com/seth16888/wxproxy/api/v1"
	"go.uber.org/zap"
)

type MPMemberRepo interface {
	FindByAppId(c context.Context, appId string) ([]*entities.MPMember, error)
	FindById(c context.Context, id string) (*entities.MPMember, error)
	UpdateRemark(c context.Context, id, remark string) error
}

type MPBlackListRepo interface {
	Block(c context.Context, appId string, openids []string) error
	Unblock(c context.Context, appId string, openids []string) error
	Query(c context.Context, appId string) ([]*entities.MPMember, error)
}

type MPMemberUsecase struct {
	repo          MPMemberRepo
	blackListRepo MPBlackListRepo
	log           *zap.Logger
	apiProxy      *APIProxyUsecase
}

func (m *MPMemberUsecase) BatchUnblock(c *gin.Context, appId string, openids []string) error {
	// TODO: 微信接口调用

	if err := m.blackListRepo.Unblock(c, appId, openids); err != nil {
		m.log.Error("unblock member error", zap.Error(err))
		return fmt.Errorf("unblock member error")
	}
	return nil
}

func (m *MPMemberUsecase) BatchBlock(c *gin.Context, appId string, openids []string) error {
	// TODO: 微信接口调用
	// app, err := GetAppInfoFromCtx(c)
	// if err != nil {
	// 	m.log.Error("get app info error", zap.Error(err))
	// 	return fmt.Errorf("get app info error")
	// }
	// mpId := app.MpId
	// accessToken, err := m.apiProxy.GetAccessToken(c, appId, mpId)
  // if err != nil {
	// 	m.log.Error("get access token error", zap.Error(err))
	// 	return fmt.Errorf("get access token error")
	// }
  // if wxErr, err := m.apiProxy.cli.GetMemberTags(c, &v1.BatchBlockRequest{
  //   AccessToken: accessToken,
  //   Openids: openids,
  // }); err != nil || wxErr.Errcode != 0 {
	// 	m.log.Error("batch block member error", zap.Error(err), zap.Any("wxErr", wxErr))
	// 	return fmt.Errorf("batch block member error")
	// }

	if err := m.blackListRepo.Block(c, appId, openids); err != nil {
		m.log.Error("block member error", zap.Error(err))
		return fmt.Errorf("block member error")
	}
	return nil
}

func (m *MPMemberUsecase) GetBlackList(c *gin.Context, appId string) ([]*entities.MPMember, error) {
	docs, err := m.blackListRepo.Query(c, appId)
	if err != nil {
		m.log.Error("query black list error", zap.Error(err))
		return nil, fmt.Errorf("query black list error")
	}
	return docs, nil
}

func (m *MPMemberUsecase) UpdateRemark(c *gin.Context, appId, id, openId, remark string) error {
	app, err := GetAppInfoFromCtx(c)
	if err != nil {
		m.log.Error("get app info error", zap.Error(err))
		return fmt.Errorf("get app info error")
	}
	mpId := app.MpId
	accessToken, err := m.apiProxy.GetAccessToken(c, appId, mpId)
  if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

  req:= v1.UpdateMemberRemarkRequest{
    AccessToken: accessToken,
    Openid: openId,
    Remark: remark,
  }
  if wxErr, err := m.apiProxy.cli.UpdateMemberRemark(c, &req); err != nil || wxErr.Errcode != 0 {
		m.log.Error("update remark error", zap.Error(err), zap.Any("wxErr", wxErr))
		return fmt.Errorf("update remark error")
	}

	_, err = m.repo.FindById(c, id)
	if err != nil {
		m.log.Error("update member remark error", zap.Error(err))
		return fmt.Errorf("data not found")
	}

	if err := m.repo.UpdateRemark(c, id, remark); err != nil {
		m.log.Error("update member remark error", zap.Error(err))
		return fmt.Errorf("update member remark error")
	}
	return nil
}

func (m *MPMemberUsecase) GetMemberInfo(c *gin.Context, appId string, id string) (*entities.MPMember, error) {
	return m.repo.FindById(c, id)
}

func (m *MPMemberUsecase) Query(c *gin.Context, appId string) ([]*entities.MPMember, error) {
	docs, err := m.repo.FindByAppId(c, appId)
	if err != nil {
		m.log.Error("query member error", zap.Error(err))
		return nil, fmt.Errorf("query member error")
	}
	return docs, nil
}

func NewMPMemberUsecase(log *zap.Logger, repo MPMemberRepo, apiProxy *APIProxyUsecase) *MPMemberUsecase {
	return &MPMemberUsecase{repo: repo, log: log, apiProxy: apiProxy}
}
