package handler

import (
	"github.com/gin-gonic/gin"
	au "github.com/seth16888/coauth/api/v1"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/internal/model/response"
	"github.com/seth16888/wxbusiness/pkg/validator"
	"go.uber.org/zap"
)

type AuthHandler struct {
	Base
	log       *zap.Logger
	ac        au.CoauthClient
	validator *validator.Validator
}

func NewAuthHandler(log *zap.Logger, validator *validator.Validator,
	ac au.CoauthClient,
) *AuthHandler {
	return &AuthHandler{log: log, validator: validator, ac: ac}
}

// GetCaptcha 获取验证码
func (h *AuthHandler) GetCaptcha(c *gin.Context) {
	ret, err := h.ac.Captcha(c, &au.CaptchaRequest{})
	if err != nil {
		h.log.Error("GetCaptcha", zap.Error(err))
		c.JSON(500, r.Error(500, "获取验证码失败"))
		return
	}

	resp := &response.GetCaptchaResp{
		CaptchaKey:   ret.CaptchaKey,
		CaptchaValue: ret.CaptchaValue,
	}
	c.JSON(200, r.SuccessData(resp))
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, r.Error(400, "参数错误"))
		return
	}
	if err := h.validator.Validate(&req); err != nil {
		c.JSON(400, r.Error(400, err.Error()))
		return
	}

	ret, err := h.ac.Login(c, &au.LoginRequest{
		Username:    req.Username,
		Password:    req.Password,
		CaptchaKey:  req.CaptchaKey,
		CaptchaCode: req.CaptchaCode,
	})
	if err != nil {
		h.log.Error("Login", zap.Error(err))
		c.JSON(400, r.Error(400, "username or password error"))
		return
	}

	resp := response.LoginResp{
		AccessToken:  ret.AccessToken,
		RefreshToken: ret.RefreshToken,
		TokenType:    ret.TokenType,
		ExpiresIn:    ret.ExpiresIn,
	}
	c.JSON(200, r.SuccessData(resp))
}
