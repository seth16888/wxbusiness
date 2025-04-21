package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/seth16888/wxbusiness/internal/data/entities"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/pkg/material"
	"github.com/seth16888/wxcommon/domain"
	"github.com/seth16888/wxcommon/hc"
	"github.com/seth16888/wxcommon/mp"
	"github.com/seth16888/wxcommon/paths"
	v1 "github.com/seth16888/wxproxy/api/v1"
	"go.uber.org/zap"
)

type MaterialRepo interface {
	Insert(c context.Context, material *entities.MPMaterial) error
	Find(c context.Context, appId string, IsPermanent bool, mediaType string,
		pageNo int64, pageSize int64) ([]*entities.MPMaterial, error)
  SaveMany(c context.Context, materials []*entities.MPMaterial) error
}

type MaterialUsecase struct {
	log      *zap.Logger
	repo     MaterialRepo
	apiProxy *APIProxyUsecase
	hc       *hc.Client
}

func (m *MaterialUsecase) Pull(c context.Context, appId string, req *request.PullMaterialReq) error {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	if req.Type == "news" {
		return m.pullNewsList(c, token, req, appId, mpId)
	}
	return m.pullMediaList(c, token, req, appId, mpId)
}

func (m *MaterialUsecase) pullNewsList(c context.Context, token string, req *request.PullMaterialReq, appId string, mpId string) error {
	params := &v1.GetMaterialListRequest{
		AccessToken: token,
		Offset:      req.Offset,
		Count:       req.Count,
		Type:        req.Type,
	}
	stream, err := m.apiProxy.cli.GetMaterialNewsList(c, params)
	if err != nil {
		m.log.Error("get material news list error", zap.Error(err))
		return fmt.Errorf("get material news list error")
	}

	var docs []*entities.MPMaterial
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			m.log.Error("receive material news list error", zap.Error(err))
			return fmt.Errorf("receive material news list error")
		}
		if resp == nil || resp.Item == nil {
			continue
		}
		// 转换为MPMaterial
		for _, item := range resp.Item {
			for _, article := range item.Articles {
				doc := &entities.MPMaterial{
					AppId:            appId,
					MpId:             mpId,
					Type:             "news",
					IsPermanent:      true,
					MediaId:          item.MediaId,
					ThumbMediaId:     article.ThumbMediaId,
					Title:            article.Title,
					Content:          article.Content,
					URL:              article.Url,
					CreatedAt:        item.UpdateTime,
					UpdatedAt:        item.UpdateTime,
					Filename:         "",
					Path:             "",
					Author:           article.Author,
					Digest:           article.Digest,
					ShowCoverPic:     article.ShowCoverPic,
					ContentSourceUrl: article.ContentSourceUrl,
				}
				docs = append(docs, doc)
			}
		}
	}
	// save to db
	if len(docs) > 0 {
		err = m.repo.SaveMany(c, docs)
		if err!= nil {
			m.log.Error("save material news list error", zap.Error(err))
			return fmt.Errorf("save material news list error")
		}
	}
	return nil
}

func (m *MaterialUsecase) pullMediaList(c context.Context, token string, req *request.PullMaterialReq, appId string, mpId string) error {
	params := &v1.GetMaterialListRequest{
		AccessToken: token,
		Offset:      req.Offset,
		Count:       req.Count,
		Type:        req.Type,
	}
	stream, err := m.apiProxy.cli.GetMaterialList(c, params)
	if err != nil {
		m.log.Error("get material list error", zap.Error(err))
		return fmt.Errorf("get material list error")
	}

	var docs []*entities.MPMaterial
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			m.log.Error("receive material news list error", zap.Error(err))
			return fmt.Errorf("receive material news list error")
		}
		if resp == nil || resp.Item == nil {
			continue
		}
		// 转换为MPMaterial
		for _, item := range resp.Item {
			doc := &entities.MPMaterial{
				AppId:       appId,
				MpId:        mpId,
				Type:        req.Type,
				IsPermanent: true,
				MediaId:     item.MediaId,
				CreatedAt:   item.UpdateTime,
				UpdatedAt:   item.UpdateTime,
				Filename:    item.Name,
				URL:         item.Url,
			}
			docs = append(docs, doc)
		}
	}
	// save to db
  if len(docs) > 0 {
		err = m.repo.SaveMany(c, docs)
		if err!= nil {
			m.log.Error("save material news list error", zap.Error(err))
			return fmt.Errorf("save material news list error")
		}
	}
	return nil
}

// UploadNewsImage 上传图文素材图片
//
// 返回: url
func (m *MaterialUsecase) UploadNewsImage(c context.Context, appId string,
	filename string, path string,
) (string, error) {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return "", fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return "", fmt.Errorf("get access token error")
	}

	url := fmt.Sprintf("https://%s%s?access_token=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Upload_News_Image,
		token,
	)
	m.log.Debug("UploadNewsImage", zap.String("url", url))

	var respBytes []byte
	respBytes, err = material.UploadFile("media", path, url, m.hc)
	if err != nil {
		m.log.Debug("UploadNewsImage", zap.Error(err))
		return "", fmt.Errorf("upload file error")
	}

	type resultT struct {
		mp.WXError
		URL string `json:"url"`
	}
	var resultVar resultT
	err = json.Unmarshal(respBytes, &resultVar)
	if err != nil {
		m.log.Debug("UploadNewsImage", zap.Error(err))
		return "", fmt.Errorf("unmarshal response error")
	}
	if resultVar.ErrCode != 0 { // business error
		m.log.Error("UploadNewsImage", zap.Any("ApiError", resultVar))
		return "", fmt.Errorf("upload file error: %d %s", resultVar.ErrCode, resultVar.ErrMsg)
	}

	now := time.Now().Unix()
	doc := entities.MPMaterial{ // 只有URL，没有MediaId
		AppId:        appId,
		MpId:         mpId,
		Type:         "ArticleImage",
		IsPermanent:  true,
		MediaId:      "",
		ThumbMediaId: "",
		URL:          resultVar.URL,
		Filename:     filename,
		Path:         path,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = m.repo.Insert(c, &doc)
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
//
// 媒体文件类型，分别有图片（image）、语音（voice）、视频（video）和缩略图（thumb）
//
// 视频素材的标题 title，不超过128个字节，超过会自动截断,
// 视频素材的描述 introduction，不超过512个字节，超过会自动截断
//
// 返回: media_id, url(只有图片素材有)
func (m *MaterialUsecase) UploadMedia(c context.Context, appId string,
	mediaType string, filename string, path string, videoTitle string, videoIntro string,
) (*mp.UploadMaterialRes, error) {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return nil, fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return nil, fmt.Errorf("get access token error")
	}

	url := fmt.Sprintf("https://%s%s?access_token=%s&type=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Upload_Media,
		token,
		mediaType,
	)
	m.log.Debug("UploadMedia", zap.String("url", url))

	var respBytes []byte
	if mediaType == "video" {
		videoForm := fmt.Sprintf(`{"title":"%s","introduction":"%s"}`, videoTitle, videoIntro)
		fields := []material.MultipartFormField{
			{
				IsFile:    true,
				Fieldname: "media",
				FilePath:  path,
				Filename:  filename,
				Value:     nil,
			},
			{
				IsFile:    false,
				Fieldname: "description",
				FilePath:  "",
				Filename:  "",
				Value:     []byte(videoForm),
			},
		}
		respBytes, err = material.PostMultipartForm(fields, url, m.hc)
		if err != nil {
			m.log.Debug("UploadMedia(video)", zap.Error(err))
			return nil, fmt.Errorf("upload video file error")
		}
	} else {
		respBytes, err = material.UploadFile("media", path, url, m.hc)
		if err != nil {
			m.log.Debug("UploadMedia", zap.Error(err))
			return nil, fmt.Errorf("upload media file error")
		}
	}
	// 返回结果
	type result struct {
		mp.WXError
		mp.UploadMaterialRes
	}
	var resultVar result
	err = json.Unmarshal(respBytes, &resultVar)
	if err != nil {
		m.log.Debug("UploadMedia", zap.Error(err))
		return nil, fmt.Errorf("unmarshal response error")
	}
	if resultVar.ErrCode != 0 { // business error
		m.log.Error("UploadMedia", zap.Any("ApiError", resultVar))
		return nil, fmt.Errorf("upload media file error: %d %s", resultVar.ErrCode, resultVar.ErrMsg)
	}

	now := time.Now().Unix()
	doc := entities.MPMaterial{
		AppId:        appId,
		MpId:         mpId,
		Type:         mediaType,
		IsPermanent:  true,
		MediaId:      resultVar.MediaID,
		ThumbMediaId: "",
		URL:          resultVar.Url,
		Filename:     filename,
		Path:         path,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	err = m.repo.Insert(c, &doc)
	if err != nil {
		m.log.Error("insert material error", zap.Error(err))
		return nil, fmt.Errorf("insert material error")
	}
	return &resultVar.UploadMaterialRes, nil
}

// UploadTemporaryMedia 上传临时素材
//
// 3天后过期
func (m *MaterialUsecase) UploadTemporaryMedia(c context.Context,
	appId string, mediaType string, filename string, path string,
) (*mp.UploadMediaRes, error) {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return nil, fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return nil, fmt.Errorf("get access token error")
	}

	url := fmt.Sprintf("https://%s%s?access_token=%s&type=%s",
		domain.GetWXAPIDomain(),
		paths.Path_Upload_Temporary_Media,
		token,
		mediaType,
	)
	m.log.Debug("UploadTemporaryMedia", zap.String("url", url))

	var respBytes []byte
	respBytes, err = material.UploadFile("media", path, url, m.hc)
	if err != nil {
		m.log.Debug("UploadTemporaryMedia", zap.Error(err))
		return nil, fmt.Errorf("upload file error")
	}
	// 返回结果
	type result struct {
		mp.WXError
		mp.UploadMediaRes
	}
	var resultVar result
	err = json.Unmarshal(respBytes, &resultVar)
	if err != nil {
		m.log.Debug("UploadTemporaryMedia", zap.Error(err))
		return nil, fmt.Errorf("unmarshal response error")
	}
	if resultVar.ErrCode != 0 { // business error
		m.log.Error("UploadTemporaryMedia", zap.Any("ApiError", resultVar))
		return nil, fmt.Errorf("upload file error: %d %s", resultVar.ErrCode, resultVar.ErrMsg)
	}

	doc := entities.MPMaterial{
		AppId:        appId,
		MpId:         mpId,
		Type:         resultVar.Type,
		IsPermanent:  false,
		MediaId:      resultVar.MediaID,
		ThumbMediaId: resultVar.ThumbMediaId,
		URL:          resultVar.URL,
		Filename:     filename,
		Path:         path,
		CreatedAt:    resultVar.CreatedAt,
		UpdatedAt:    resultVar.CreatedAt,
	}
	err = m.repo.Insert(c, &doc)
	if err != nil {
		m.log.Error("insert material error", zap.Error(err))
		return nil, fmt.Errorf("insert material error")
	}

	return &resultVar.UploadMediaRes, nil
}

// DeleteMaterial 删除素材(永久)
func (m *MaterialUsecase) DeleteMaterial(c context.Context, appId string, mediaId string) error {
	mpIdVar := c.Value("MP_ID")
	if mpIdVar == nil {
		m.log.Error("get mp id error")
		return fmt.Errorf("get mp id error")
	}
	mpId := mpIdVar.(string)

	token, err := m.apiProxy.GetAccessToken(c, appId, mpId)
	if err != nil {
		m.log.Error("get access token error", zap.Error(err))
		return fmt.Errorf("get access token error")
	}

	params := &v1.DeleteMaterialReq{
		AccessToken: token,
		MediaId:     mediaId,
	}
	wxErr, err := m.apiProxy.cli.DeleteMaterial(c, params)
	if err != nil {
		m.log.Error("DeleteMaterial error", zap.Error(err))
		return fmt.Errorf("DeleteMaterial error")
	}
	if wxErr.Errcode != 0 {
		m.log.Error("DeleteMaterial", zap.Any("ApiError", wxErr))
		return fmt.Errorf("DeleteMaterial error: %d %s", wxErr.Errcode, wxErr.Errmsg)
	}
	return nil
}

func NewMaterialUsecase(log *zap.Logger, repo MaterialRepo,
	apiProxy *APIProxyUsecase, hc *hc.Client,
) *MaterialUsecase {
	return &MaterialUsecase{log: log, repo: repo, apiProxy: apiProxy, hc: hc}
}
