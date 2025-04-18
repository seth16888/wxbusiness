package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/pkg/validator"
	"go.uber.org/zap"
)

// QRCodeHandler 二维码处理器
type QRCodeHandler struct {
  Base
  log *zap.Logger
	uc  *biz.MpQRCodeUsecase
  validator *validator.Validator
}

// NewQRCodeHandler
func NewQRCodeHandler(log *zap.Logger, uc *biz.MpQRCodeUsecase,
  validator *validator.Validator) *QRCodeHandler {
  return &QRCodeHandler{log:log, uc:uc, validator: validator}
}

// CreateTemporary
func (h *QRCodeHandler) CreateTemporary(ctx *gin.Context) {
  // 路径参数
  appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
  // 请求参数
  var req request.CreateQRCodeReq
  if err := h.BindAndValidate(ctx, h.validator, &req);err!=nil {
    ctx.JSON(400, r.Error(400, err.Error()))
    return
  }

  c := ctx
  res,err := h.uc.CreateTemporary(c, appId, &req)
  if err != nil {
    ctx.JSON(400, r.Error(400, err.Error()))
    return
  }

  ctx.JSON(200, r.SuccessData(res))
}

// CreateLimit
func (h *QRCodeHandler) CreateLimit(ctx *gin.Context) {
  // 路径参数
  appId, err := h.GetPID(ctx)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
  // 请求参数
  var req request.CreateQRCodeReq
  if err := h.BindAndValidate(ctx, h.validator, &req);err!=nil {
    ctx.JSON(400, r.Error(400, err.Error()))
    return
  }

  c := ctx
  res,err := h.uc.CreateLimit(c, appId, &req)
  if err != nil {
    ctx.JSON(400, r.Error(400, err.Error()))
    return
  }

  ctx.JSON(200, r.SuccessData(res))
}

// GetURL
func (h *QRCodeHandler) GetURL(ctx *gin.Context) {
  // 请求参数
  var req request.CreateQRCodeReq
  if err := h.BindAndValidate(ctx, h.validator, &req);err!=nil {
    ctx.JSON(400, r.Error(400, err.Error()))
    return
  }
  // query 参数
  ticket := ctx.Query("ticket")
  if ticket == "" {
    ctx.JSON(400, r.Error(400, "ticket is required"))
    return
  }

  c := ctx
  url := h.uc.GetURL(c, ticket)

  ctx.JSON(200, r.SuccessData(url))
}
