package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/pkg/validator"
	"go.uber.org/zap"
)

type PlatformAppHandler struct {
	Base
	validator *validator.Validator
	log       *zap.Logger
	uc        *biz.PlatformAppUsecase
}

func NewPlatformAppHandler(
	log *zap.Logger,
	validator *validator.Validator,
	uc *biz.PlatformAppUsecase,
) *PlatformAppHandler {
	return &PlatformAppHandler{
		log:       log,
		validator: validator,
		uc:        uc,
	}
}

// Create
func (h *PlatformAppHandler) Create(ctx *gin.Context) {
	var req request.CreateAppReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, r.Error[any](400, err.Error()))
		return
	}
	if err := h.validator.Validate(&req); err != nil {
		ctx.JSON(400, r.Error[any](400, err.Error()))
		return
	}
	// User Id
	userId, err := h.GetUserId(ctx)
	if err != nil {
		ctx.JSON(401, r.Error[any](401, "unauthorized"))
		return
	}

	res, err := h.uc.Create(ctx, userId, &req)
	if err != nil {
		ctx.JSON(400, r.Error[any](400, err.Error()))
		return
	}

	ctx.JSON(200, r.Success(res))
}
