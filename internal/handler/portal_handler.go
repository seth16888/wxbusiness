package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/message"
	"github.com/seth16888/wxbusiness/internal/model/r"
	"github.com/seth16888/wxbusiness/internal/model/request"
	"github.com/seth16888/wxbusiness/pkg/logger"
	"github.com/seth16888/wxbusiness/pkg/validator"
	commHelpers "github.com/seth16888/wxcommon/helpers"
	"go.uber.org/zap"
)

type PortalHandler struct {
	Base
	validator *validator.Validator
	log       *zap.Logger
	uc        *biz.PortalUsecase
}

func NewPortalHandler(log *zap.Logger, validator *validator.Validator,
	uc *biz.PortalUsecase,
) *PortalHandler {
	return &PortalHandler{log: log, validator: validator, uc: uc}
}

// Verify
func (h *PortalHandler) Verify(ctx *gin.Context) {
	var params request.VerifyPortalReq
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}
	// App Id
	params.AppID = ctx.Param("id")
	if params.AppID == "" {
		ctx.JSON(400, r.Error(400, "invalid appid"))
		return
	}

	logger.Debugf("wx verify req,appid: %s,timestamp: %s,nonce: %s,signature: %s",
		params.AppID, params.Timestamp, params.Nonce, params.Signature)

	if err := h.validator.Validate(&params); err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}

	if _, err := h.uc.Verify(ctx, params.AppID, params.Timestamp,
		params.Nonce, params.Signature, params.EchoStr); err != nil {
		ctx.JSON(400, r.Error(400, "invalid signature"))
		return
	}

	// 返回EchoStr 字符串
	ctx.String(200, params.EchoStr)
}

// Portal
func (h *PortalHandler) Portal(ctx *gin.Context) {
	// App Id
	appId := ctx.Param("id")
	if len(appId) == 0 {
		ctx.JSON(400, r.Error(400, "invalid mpId"))
		return
	}
	// 获取mp token
	mpAPP, err := h.uc.GetMpApp(ctx, appId)
	if err != nil {
		logger.Errorf("get mp token error: %v", err)
		ctx.JSON(400, r.Error(400, "get mp token error"))
		return
	}

	var params request.PortalMessageReq
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(400, r.Error(400, err.Error()))
		return
	}

	if len(params.EncryptType) == 0 {
    // 明文模式
		params.EncryptType = "plain"
    // 验证签名
    if !commHelpers.ValidatePortalReq(params.Timestamp, params.Nonce, params.Signature, mpAPP.Token) {
      logger.Errorf("signature error")
      ctx.JSON(400, r.Error(400, "invalid signature"))
      return
    }
	} else {
    // 安全模式
    var encryptMsg request.EncryptMessageReq
    if err := ctx.ShouldBind(&encryptMsg); err != nil {
      logger.Errorf("bind encrypt message error: %v", err)
      ctx.JSON(400, r.Error(400, "invalid encrypt message"))
      return
    }
    // 验证签名
    if !commHelpers.ValidateEncryptSignature(mpAPP.Token, params.Timestamp, params.Nonce, params.MsgSignature, encryptMsg.Encrypt) {
      logger.Errorf("signature error")
      ctx.JSON(400, r.Error(400, "invalid signature"))
      return
    }
    // 解密消息
    cipherText := []byte(encryptMsg.Encrypt)
    aesKey := []byte(mpAPP.EncodingAesKey[16:32])
    iv := []byte(mpAPP.EncodingAesKey[0:16])

    // 验证 cipherText 的长度是否是 16 字节的整数倍
    blockSize := 16
    if len(cipherText)%blockSize != 0 {
      h.log.Error("cipherText length is not a multiple of block size", zap.Int("length", len(cipherText)), zap.Int("blockSize", blockSize))
      ctx.JSON(400, r.Error(400, "invalid cipherText length"))
      return
    }

    rawData, err := commHelpers.AESDecryptData(cipherText, aesKey, iv)
    if err != nil {
      logger.Errorf("decrypt message error: %v", err)
      ctx.JSON(400, r.Error(400, "decrypt message error"))
      return
    }
    logger.Debugf("raw data: %s", string(rawData))
    // TODO: 解析加密消息
    ctx.String(200, "success")
		return
  }
	logger.Debugf("wx push: %+v", params)

	// 绑定body
	msgDomain := message.NewMessageDomain(ctx.Request.Body)
	msgDomain.SetContentType(ctx.Request.Header.Get("Content-Type"))
	msgDomain.SetEncryptType(params.EncryptType)
	msgDomain.OpenId = params.OpenID
	msgDomain.Timestamp = params.Timestamp
	msgDomain.Nonce = params.Nonce
	// 解析消息
	err = msgDomain.Unmarshal()
	if err != nil {
		logger.Errorf("unmarshal error: %v", err)
		ctx.JSON(400, r.Error(400, "invalid xml"))
		return
	}

	resultChan := message.MessageWorker(msgDomain)
	select {
	case result := <-resultChan:
		if len(result) == 0 {
			ctx.String(200, "success")
			return
		} else {
			logger.Debugf("reply: %s", string(result))
			ctx.Writer.Header().Set("Content-Type", "text/xml;charset=utf-8")
			ctx.Writer.WriteHeader(200)
			ctx.Writer.Write(result)
			return
		}
	case <-time.After(time.Second * 4):
		logger.Errorf("wx portal reply timeout")
		ctx.String(200, "success")
		return
	}
}
