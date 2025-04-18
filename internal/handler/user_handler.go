package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxbusiness/pkg/validator"
	"go.uber.org/zap"
)

type UserHandler struct {
	log       *zap.Logger
	validator *validator.Validator
	uc        *biz.UserUsecase
}

func NewUserHandler(log *zap.Logger, validator *validator.Validator,
  uc *biz.UserUsecase) *UserHandler {
	return &UserHandler{log: log, validator: validator, uc: uc}
}

// ListMPApps
// @Summary 获取用户公众号列表
func (u *UserHandler) ListMPApps(c *gin.Context) {
	userId := c.Param("userId")
	if userId == "" {
		c.JSON(400, r.Error(400, "userId is required"))
		return
	}

	res, err := u.uc.ListMPApps(c, userId)
	if err != nil {
		c.JSON(500, r.Error(500, err.Error()))
		return
	}

	c.JSON(200, r.SuccessData(res))
}
