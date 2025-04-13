package request

// VerifyPortalReq 开发接口验证请求
type VerifyPortalReq struct {
	AppID     string `form:"appId"`
	Signature string `form:"signature" binding:"required"`
	Timestamp string `form:"timestamp" binding:"required"`
	Nonce     string `form:"nonce" binding:"required"`
	EchoStr   string `form:"echostr" binding:"required"`
}

// WXPushMessage 推送消息请求
type WXPushMessage struct {
	AppID       string `json:"appId" form:"appId"`
	EncryptType string `json:"encryptType" form:"encrypt_type"`
	Nonce       string `json:"nonce" form:"nonce"`
	Timestamp   string `json:"timestamp" form:"timestamp"`
	OpenID      string `json:"openid" form:"openid"`
}

// CreateAppReq 创建应用请求
type CreateAppReq struct {
	Name         string `json:"name" binding:"required"`
	PicUrl       string `json:"picUrl" binding:"required"`
	Introduction string `json:"introduction" binding:"required"`
	AppId        string `json:"appId" binding:"required"`
	AppSecret    string `json:"appSecret" binding:"required"`
	Type         int64  `json:"type" binding:"required"`
}

