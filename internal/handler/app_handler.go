package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/pkg/validator"
	"go.uber.org/zap"
)

type AppHandler struct {
	Base
	validator *validator.Validator
	log       *zap.Logger
	uc        *biz.AppUsecase
}

func NewAppHandler(log *zap.Logger, validator *validator.Validator,
  uc *biz.AppUsecase) *AppHandler {
	return &AppHandler{log: log, validator: validator, uc: uc}
}

// Create
func (h *AppHandler) Create(ctx *gin.Context) {
	var req request.CreateAppReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	if err := h.validator.Validate(&req); err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	// User Id
	userId, err := h.GetUserId(ctx)
	if err != nil {
		ctx.JSON(401, r.Error(401, "unauthorized"))
		return
	}

	res, err := h.uc.Create(ctx, userId, &req)
	if err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}

	ctx.JSON(200, r.SuccessData(res))
}

// GetById
func (h *AppHandler) GetById(ctx *gin.Context) {
	appId := ctx.Param("id")
	if appId == "" {
		ctx.JSON(400, r.Error(400, "appId is required"))
		return
	}

	res, err := h.uc.GetById(ctx, appId)
	if err != nil {
		ctx.JSON(500, r.Error(500, err.Error()))
		return
	}

	ctx.JSON(200, r.SuccessData(res))
}
