package handler

import (
	"path"
	"path/filepath"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxcommon/helpers"
	"go.uber.org/zap"
)

type MaterialHandler struct {
	Base
	log *zap.Logger
	uc  *biz.MaterialUsecase
}

func NewMaterialHandler(log *zap.Logger, uc *biz.MaterialUsecase) *MaterialHandler {
	return &MaterialHandler{log: log, uc: uc}
}

// UploadTemporaryMedia
func (h *MaterialHandler) UploadTemporaryMedia(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	mediaType := ctx.Query("type")
	// 允许的素材类型
	allowedTypes := []string{"image", "voice", "video", "thumb"}
	if !slices.Contains(allowedTypes, mediaType) {
		ctx.JSON(400, r.Error(400, "type allowed types are image, voice, video, thumb"))
		return
	}
	// 上传文件
	file, err := ctx.FormFile("file")
	if err != nil || file == nil {
		ctx.JSON(400, r.Error(400, "file not found"))
		return
	}
	// 检查文件类型、大小
	ext := filepath.Ext(file.Filename)
	switch mediaType {
	case "image":
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			ctx.JSON(400, r.Error(400, "image file type not allowed"))
			return
		}
		if file.Size > 10*1024*1024 {
			ctx.JSON(400, r.Error(400, "image file size too large"))
			return
		}
	case "voice":
		if ext != ".mp3" && ext != ".amr" {
			ctx.JSON(400, r.Error(400, "voice file type not allowed"))
			return
		}
		if file.Size > 2*1024*1024 {
			ctx.JSON(400, r.Error(400, "voice file size too large"))
			return
		}
	case "video":
		if ext != ".mp4" {
			ctx.JSON(400, r.Error(400, "video file type not allowed"))
			return
		}
		if file.Size > 10*1024*1024 {
			ctx.JSON(400, r.Error(400, "video file size too large"))
			return
		}
	case "thumb":
		if ext != ".jpg" {
			ctx.JSON(400, r.Error(400, "thumb file type not allowed"))
			return
		}
		if file.Size > 64*1024 {
			ctx.JSON(400, r.Error(400, "thumb file size too large"))
			return
		}
	}

	// 保存文件
  filename := helpers.UUID() + ext
	dst := path.Join("uploads", filename)
	if err = ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.JSON(500, r.Error(500, "save file error"))
		return
	}

  // 上传到微信
  c := ctx
  res, err := h.uc.UploadTemporaryMedia(c, appId, mediaType, filename, dst)
  if err!= nil {
    ctx.JSON(500, r.Error(500, err.Error()))
    return
  }

  ctx.JSON(200, r.SuccessData(res))
}

// UploadNewsImage 上传图文消息内的图片获取URL
//
// 本接口所上传的图片不占用公众号的素材库中图片数量的100000个的限制。
// 图片仅支持jpg/png格式，大小必须在1MB以下。
func (h *MaterialHandler) UploadNewsImage(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err!= nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	// 上传文件
	file, err := ctx.FormFile("file")
	if err!= nil || file == nil {
		ctx.JSON(400, r.Error(400, "file not found"))
		return
	}
  // 检查文件类型、大小
	ext := filepath.Ext(file.Filename)
	if ext!= ".jpg" && ext!= ".png" {
		ctx.JSON(400, r.Error(400, "image file type not allowed"))
		return
	}
	if file.Size > 1*1024*1024 {
		ctx.JSON(400, r.Error(400, "image file size too large"))
		return
	}
  // 保存文件
  filename := helpers.UUID() + ext
	dst := path.Join("uploads", filename)
	if err = ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.JSON(500, r.Error(500, "save file error"))
		return
	}
  // 上传到微信
  c := ctx
  res, err := h.uc.UploadNewsImage(c, appId, filename, dst)
  if err!= nil {
    ctx.JSON(500, r.Error(500, err.Error()))
    return
  }
  ctx.JSON(200, r.SuccessData(res))
}

// UploadMedia 上传其他类型永久素材
// 1、最近更新：永久图片素材新增后，将带有URL返回给开发者，开发者可以在腾讯系域名内使用（腾讯系域名外使用，图片将被屏蔽）。
// 2、公众号的素材库保存总数量有上限：图文消息素材、图片素材上限为100000，其他类型为1000。
// 3、素材的格式大小等要求与公众平台官网一致：
//
//	图片（image）: 10M，支持bmp/png/jpeg/jpg/gif格式
//	语音（voice）：2M，播放长度不超过60s，mp3/wma/wav/amr格式
//	视频（video）：10MB，支持MP4格式
//	缩略图（thumb）：64KB，支持JPG格式
func (h *MaterialHandler) UploadMedia(ctx *gin.Context) {
	// 路径参数
	appId, err := h.GetPID(ctx)
	if err!= nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
  // query 参数
  mediaType := ctx.Query("type")
  // 允许的素材类型
  allowedTypes := []string{"image", "voice", "video", "thumb"}
  if!slices.Contains(allowedTypes, mediaType) {
    ctx.JSON(400, r.Error(400, "type allowed types are image, voice, video, thumb"))
    return
  }
	// 上传文件
	file, err := ctx.FormFile("file")
	if err!= nil || file == nil {
		ctx.JSON(400, r.Error(400, "file not found"))
		return
	}
  // 检查文件类型、大小
	ext := filepath.Ext(file.Filename)
	switch mediaType {
	case "image":
		if ext != ".bmp" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			ctx.JSON(400, r.Error(400, "image file type not allowed"))
			return
		}
		if file.Size > 10*1024*1024 {
			ctx.JSON(400, r.Error(400, "image file size too large"))
			return
		}
	case "voice":
		if ext != ".mp3" && ext != ".amr" && ext!= ".wma" && ext!= ".wav" {
			ctx.JSON(400, r.Error(400, "voice file type not allowed"))
			return
		}
		if file.Size > 2*1024*1024 {
			ctx.JSON(400, r.Error(400, "voice file size too large"))
			return
		}
	case "video":
		if ext != ".mp4" {
			ctx.JSON(400, r.Error(400, "video file type not allowed"))
			return
		}
		if file.Size > 10*1024*1024 {
			ctx.JSON(400, r.Error(400, "video file size too large"))
			return
		}
	case "thumb":
		if ext != ".jpg" {
			ctx.JSON(400, r.Error(400, "thumb file type not allowed"))
			return
		}
		if file.Size > 64*1024 {
			ctx.JSON(400, r.Error(400, "thumb file size too large"))
			return
		}
	}
  // 如果是video，需要额外检查title和description
  videoTitle := ""
  videoIntro := ""
  if mediaType == "video" {
    videoTitle = ctx.Query("title")
    if videoTitle == "" {
      ctx.JSON(400, r.Error(400, "video title not found"))
      return
    }
    videoIntro = ctx.Query("introduction")
  }

  // 保存文件
  filename := helpers.UUID() + ext
	dst := path.Join("uploads", filename)
	if err = ctx.SaveUploadedFile(file, dst); err!= nil {
		ctx.JSON(500, r.Error(500, "save file error"))
		return
	}

  // 上传到微信
  c := ctx
  res, err := h.uc.UploadMedia(c, appId, mediaType, filename, dst,
    videoTitle, videoIntro)
  if err!= nil {
    ctx.JSON(500, r.Error(500, err.Error()))
    return
  }

  ctx.JSON(200, r.SuccessData(res))
}

// GetMaterialList 获取素材列表
//
func (h *MaterialHandler) GetMaterialList(ctx *gin.Context) {
  // 路径参数
	appId, err := h.GetPID(ctx)
	if err!= nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
  // query 参数
  var pageNo int64 = 1
  var pageSize int64 = 10
  var pagingQ request.PagingQuery
  if err := ctx.ShouldBindQuery(&pagingQ); err == nil {
    pageNo = pagingQ.PageNo
    pageSize = pagingQ.PageSize
  }
  mediaType := ctx.Query("type")  // 为空-查询全部类型
  IsPermanent := false
  permanent := ctx.Query("permanent")
  if len(permanent) >0 {
    IsPermanent = true
  }

  c:=ctx
  res,err:= h.uc.GetMaterialList(c, appId, IsPermanent, mediaType, pageNo, pageSize)
  if err!= nil {
    ctx.JSON(500, r.Error(500, err.Error()))
    return
  }
  ctx.JSON(200, r.SuccessData(res))
}
