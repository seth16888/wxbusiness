package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/seth16888/wxbusiness/internal/data/entities"
	v1 "github.com/seth16888/wxproxy/api/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type MemberTagRepo interface {
	Create(ctx context.Context, tag *entities.MemberTag) error
	Delete(ctx context.Context, appId string, tagId int64) error
	GetByTagId(ctx context.Context, appId string, tagId int64) (*entities.MemberTag, error)
	Update(ctx context.Context, id primitive.ObjectID, tag *entities.MemberTag) error
	Query(ctx context.Context, appId string) ([]*entities.MemberTag, error)
}

// MemberTagUsecase 会员标签
type MemberTagUsecase struct {
	repo     MemberTagRepo
	log      *zap.Logger
	apiProxy *APIProxyUsecase
}

func NewMemberTagUsecase(repo MemberTagRepo, log *zap.Logger,
	apiProxy *APIProxyUsecase,
) *MemberTagUsecase {
	return &MemberTagUsecase{repo: repo, log: log, apiProxy: apiProxy}
}

func (u *MemberTagUsecase) Create(ctx context.Context, appId, tagName string) error {
	app, err := GetAppInfoFromCtx(ctx)
	if err != nil {
		u.log.Error("get app info error", zap.Error(err))
		return fmt.Errorf("get app info error")
	}
	// appId := app.ID.Hex()
	mpId := app.MpId

	token, err := u.apiProxy.GetAccessToken(ctx, appId, mpId)
  if err != nil {
		u.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	params := v1.CreateTagRequest{
		Name:        tagName,
		AccessToken: token,
	}
	reply, err := u.apiProxy.cli.CreateTag(ctx, &params)
	if err != nil {
		u.log.Error("create tag error", zap.Error(err))
		return fmt.Errorf("create tag error")
	}

	now := time.Now().Unix()
	tag := &entities.MemberTag{
		AppId:     appId,
		MpId:      mpId,
		Name:      tagName,
		TagId:     reply.Tag.Id,
		Count:     reply.Tag.Count,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Create(ctx, tag); err != nil {
		u.log.Error("create tag error", zap.Error(err))
		return fmt.Errorf("create tag error")
	}

	return nil
}

// Query
func (u *MemberTagUsecase) Query(ctx context.Context, appId string) ([]*entities.MemberTag, error) {
	// 不需要查询微信后台数据，直接从数据库查询
	tags, err := u.repo.Query(ctx, appId)
	if err != nil {
		u.log.Error("query tag error", zap.Error(err))
		return nil, fmt.Errorf("query tag error")
	}

	return tags, nil
}

// Update
func (u *MemberTagUsecase) Update(ctx context.Context, appId,tagName string  ,tagId int64 ) error {
	mpIdVar := ctx.Value("MP_ID")
  if mpIdVar == nil {
		u.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
  mpId := mpIdVar.(string)
	token, err := u.apiProxy.GetAccessToken(ctx, appId, mpId)
  if err != nil {
		u.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

  params := v1.UpdateTagRequest{
		Id:      tagId,
		Name:       tagName,
		AccessToken: token,
	}
  wxErr, err := u.apiProxy.cli.UpdateTag(ctx, &params)
  if err != nil || wxErr.Errcode != 0 {
		u.log.Error("update tag error", zap.Error(err), zap.Any("wxErr", wxErr))
		return fmt.Errorf("update tag error")
	}

	// 获取标签信息
	tag, err := u.repo.GetByTagId(ctx, appId, tagId)
	if err != nil {
		u.log.Error("get tag error", zap.Error(err))
		return fmt.Errorf("get tag error")
	}

	// 更新标签信息
	tag.Name = tagName
	if err := u.repo.Update(ctx, tag.ID, tag); err != nil {
		u.log.Error("update tag error", zap.Error(err))
		return fmt.Errorf("update tag error")
	}

	return nil
}

// Delete
func (u *MemberTagUsecase) Delete(ctx context.Context, appId string, tagId int64) error {
	mpIdVar := ctx.Value("MP_ID")
  if mpIdVar == nil {
		u.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
  mpId := mpIdVar.(string)
	token, err := u.apiProxy.GetAccessToken(ctx, appId, mpId)
  if err != nil {
		u.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

  params := v1.DeleteTagRequest{
		Id:      tagId,
		AccessToken: token,
	}
  wxErr, err := u.apiProxy.cli.DeleteTag(ctx, &params)
  if err != nil || wxErr.Errcode != 0 {
		u.log.Error("delete tag error", zap.Error(err), zap.Any("wxErr", wxErr))
		return fmt.Errorf("delete tag error")
	}

	err = u.repo.Delete(ctx, appId, tagId)
	if err != nil {
		u.log.Error("delete tag error", zap.Error(err))
		return fmt.Errorf("delete tag error")
	}
	return nil
}
