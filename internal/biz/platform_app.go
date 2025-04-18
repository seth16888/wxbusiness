package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/seth16888/wxbusiness/internal/data/entities"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/pkg/helpers"
	"go.uber.org/zap"
)

type AppRepo interface {
	Get(ctx context.Context, id string) (*entities.PlatformApp, error)
	Create(ctx context.Context, app *entities.PlatformApp) (string, error)
	// Update(ctx context.Context, app *entities.PlatformApp) error
	// Delete(ctx context.Context, id string) error
	// List(ctx context.Context, userId uint64) ([]*entities.PlatformApp, error)
	UpdateStatus(ctx context.Context, id string, status int) error
  GetByMPId(ctx context.Context, appId string) (*entities.PlatformApp, error)
}

type AppUsecase struct {
	repo AppRepo
	log  *zap.Logger
}

func NewAppUsecase(repo AppRepo, logger *zap.Logger) *AppUsecase {
	return &AppUsecase{
		repo: repo,
		log:  logger,
	}
}

func (p *AppUsecase) Create(ctx context.Context,
	userId string, req *request.CreateAppReq,
) (string, error) {
	now := time.Now().Unix()
	token := helpers.RandomString(8)
	encodingAESKey := helpers.RandomString(43)
	// 加密方式
	encodingType := 1

	app := &entities.PlatformApp{
		UserId:         userId,
		Name:           req.Name,
		Type:           req.Type,
		Token:          token,
		EncodingAesKey: encodingAESKey,
		EncodingType:   encodingType,
		MpId:          req.AppId,
		AppSecret:      req.AppSecret,
		Status:         0,
		Introduction:   req.Introduction,
		PicUrl:         req.PicUrl,
		CreatedAt:      now,
		UpdatedAt:      now,
		DeletedAt:      0,
		Version:        0,
	}

	id, err := p.repo.Create(ctx, app)
	if err != nil {
		p.log.Error("createApp", zap.Error(err))
		return "", fmt.Errorf("保存数据失败")
	}

	return id, nil
}

// GetById
func (p *AppUsecase) GetById(ctx context.Context, appId string) (*entities.PlatformApp, error) {
	entity,err:= p.repo.Get(ctx, appId)
	if err!= nil {
		p.log.Error("getAppById", zap.Error(err))
		return nil, fmt.Errorf("获取数据失败")
	}

  // TODO: 缓存
	return entity, nil
}
