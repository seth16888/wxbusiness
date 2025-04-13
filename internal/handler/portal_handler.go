package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/pkg/logger"
	"github.com/seth16888/wxbusiness/pkg/validator"
	"go.uber.org/zap"
)

type PortalHandler struct {
	Base
	validator *validator.Validator
	log       *zap.Logger
	uc        *biz.PortalUsecase
}

func NewPortalHandler(
	log *zap.Logger,
	validator *validator.Validator,
	uc *biz.PortalUsecase,
) *PortalHandler {
	return &PortalHandler{
		log:       log,
		validator: validator,
		uc:        uc,
	}
}

// Verify
func (h *PortalHandler) Verify(ctx *gin.Context) {
	var params request.VerifyPortalReq
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(400, r.Error[any](400, err.Error()))
		return
	}
  // App Id
  params.AppID = ctx.Param("id")
  if params.AppID == "" {
    ctx.JSON(400, r.Error[any](400, "invalid appid"))
    return
  }

	logger.Debugf("wx verify req,appid: %s,timestamp: %s,nonce: %s,signature: %s",
		params.AppID, params.Timestamp, params.Nonce, params.Signature)

	if err := h.validator.Validate(&params); err != nil {
		ctx.JSON(400, r.Error[any](400, err.Error()))
		return
	}

	if _, err := h.uc.Verify(ctx, params.AppID, params.Timestamp,
    params.Nonce, params.Signature, params.EchoStr); err != nil {
		ctx.JSON(400, r.Error[any](400, "invalid signature"))
		return
	}

	// 返回EchoStr 字符串
	ctx.String(200, params.EchoStr)
}

// Portal
func (h *PortalHandler) Portal(ctx *gin.Context) {
}
