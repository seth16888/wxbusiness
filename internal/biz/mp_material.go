package biz

import (
	"context"
	"fmt"
	"time"

	"github.com/seth16888/wxbusiness/internal/data/entities"
	"github.com/seth16888/wxcommon/mp"
	"go.uber.org/zap"
)

type MaterialRepo interface {
	Insert(c context.Context, material *entities.MPMaterial) error
	Find(c context.Context, appId string, IsPermanent bool, mediaType string,
		pageNo int64, pageSize int64) ([]*entities.MPMaterial, error)
}

type MaterialUsecase struct {
	log  *zap.Logger
	repo MaterialRepo
}

// UploadNewsImage 上传图文素材图片
func (m *MaterialUsecase) UploadNewsImage(c context.Context, appId string,
	filename string, path string,
) (string, error) {
	// TODO: 调用微信上传图片接口

	url := ""
	now := time.Now().Unix()
	doc := entities.MPMaterial{ // 只有URL，没有MediaId
		AppId:        appId,
		MpId:         "",
		Type:         "NewsImage",
		IsPermanent:  true,
		MediaId:      "",
		ThumbMediaId: "",
		URL:          url,
		Filename:     filename,
		Path:         path,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err := m.repo.Insert(c, &doc)
	if err != nil {
		m.log.Error("insert material error", zap.Error(err))
		return "", fmt.Errorf("insert material error")
	}
	return url, nil
}

// GetMaterialList 返回素材列表-永久素材
func (m *MaterialUsecase) GetMaterialList(c context.Context, appId string,
	IsPermanent bool, mediaType string, pageNo int64, pageSize int64,
) ([]*entities.MPMaterial, error) {
	return m.repo.Find(c, appId, IsPermanent, mediaType, pageNo, pageSize)
}

// UploadMedia 上传素材 - 永久素材
func (m *MaterialUsecase) UploadMedia(c context.Context, appId string,
	mediaType string, filename string, path string, videoTitle string, videoIntro string,
) (*mp.UploadMaterialRes, error) {
	// TODO: 调用微信上传图片接口

	url := ""
	mediaId := ""
	now := time.Now().Unix()
	doc := entities.MPMaterial{
		AppId:        appId,
		MpId:         "",
		Type:         mediaType,
		IsPermanent:  true,
		MediaId:      mediaId,
		ThumbMediaId: "",
		URL:          url,
		Filename:     filename,
		Path:         path,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err := m.repo.Insert(c, &doc)
	if err != nil {
		m.log.Error("insert material error", zap.Error(err))
		return nil, fmt.Errorf("insert material error")
	}
	return &mp.UploadMaterialRes{
		MediaID: mediaId,
		Url:     url,
	}, nil
}

// UploadTemporaryMedia 上传临时素材
//
// 3天后过期
func (m *MaterialUsecase) UploadTemporaryMedia(c context.Context,
	appId string, mediaType string, filename string, path string,
) (*mp.UploadMediaRes, error) {
	// TODO: 调用微信上传图片接口
	mediaId := ""
	url := ""
	thumbMediaId := ""
	createdAt := time.Now().Unix()

	doc := entities.MPMaterial{
		AppId:        appId,
		MpId:         "",
		Type:         mediaType,
		IsPermanent:  false,
		MediaId:      mediaId,
		ThumbMediaId: thumbMediaId,
		URL:          url,
		Filename:     filename,
		Path:         path,
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}
	err := m.repo.Insert(c, &doc)
	if err != nil {
		m.log.Error("insert material error", zap.Error(err))
		return nil, fmt.Errorf("insert material error")
	}
	return &mp.UploadMediaRes{
		Type:         mediaType,
		MediaID:      mediaId,
		URL:          url,
		ThumbMediaId: thumbMediaId,
		CreatedAt:    createdAt,
	}, nil
}

func NewMaterialUsecase(log *zap.Logger, repo MaterialRepo) *MaterialUsecase {
	return &MaterialUsecase{log: log, repo: repo}
}
