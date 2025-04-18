package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"go.uber.org/zap"
)

type MemberTagHandler struct {
  Base
	log *zap.Logger
	uc  *biz.MemberTagUsecase
}

// NewMemberTagHandler creates a new MemberTagHandler.
func NewMemberTagHandler(log *zap.Logger,uc *biz.MemberTagUsecase) *MemberTagHandler {
	return &MemberTagHandler{log: log, uc: uc}
}

// Create
func (h *MemberTagHandler) Create(ctx *gin.Context) {
  // 路径参数
  appId, err := h.GetPID(ctx)
  // query 参数
  tagName := ctx.Query("name")
  if err != nil || tagName == "" {
    ctx.JSON(400, r.Error(400, "参数错误"))
    return
  }

  c:= ctx
  if err := h.uc.Create(c, appId, tagName); err!= nil {
    h.log.Error("create tag error", zap.Error(err))
    ctx.JSON(500, r.Error(500, "创建标签失败"))
    return
  }
  ctx.JSON(200, r.Success())
}

// Delete
func (h *MemberTagHandler) Delete(ctx *gin.Context) {
  // 路径参数
  appId,err := h.GetPID(ctx)
  tagIdVar := ctx.Param("tagId")
  if err !=nil || tagIdVar == "" {
    ctx.JSON(400, r.Error(400, "参数错误"))
    return
  }
  tagId, err := strconv.ParseInt(tagIdVar, 10, 64)
  if err !=nil || tagId == 0 {
    ctx.JSON(400, r.Error(400, "tagId参数错误"))
    return
  }

  c:= ctx
  if err := h.uc.Delete(c, appId, tagId); err!= nil {
    h.log.Error("delete tag error", zap.Error(err))
    ctx.JSON(500, r.Error(500, "删除标签失败"))
    return
  }
  ctx.JSON(200, r.Success())
}

// Update
func (h *MemberTagHandler) Update(ctx *gin.Context) {
  // 路径参数
  appId,err := h.GetPID(ctx)
  tagIdVar := ctx.Param("tagId")
  // query 参数
  tagName := ctx.Query("name")
  if err !=nil || tagIdVar == "" || tagName == "" {
    ctx.JSON(400, r.Error(400, "参数错误"))
    return
  }
  tagId, err := strconv.ParseInt(tagIdVar, 10, 64)
  if err !=nil || tagId == 0 {
    ctx.JSON(400, r.Error(400, "tagId参数错误"))
    return
  }

  c:= ctx
  if err := h.uc.Update(c, appId, tagName, tagId); err!= nil {
    h.log.Error("update tag error", zap.Error(err))
    ctx.JSON(500, r.Error(500, "更新标签失败"))
    return
  }
  ctx.JSON(200, r.Success())
}

// Query
func (h *MemberTagHandler) Query(ctx *gin.Context) {
  // 路径参数
  appId,err := h.GetPID(ctx)
  if err != nil {
    ctx.JSON(400, r.Error(400, "参数错误"))
    return
  }
  c:= ctx
  tags, err := h.uc.Query(c, appId)
  if err!= nil {
    h.log.Error("query tag error", zap.Error(err))
    ctx.JSON(500, r.Error(500, "查询标签失败"))
    return
  }

  ctx.JSON(200, r.SuccessData(tags))
}
