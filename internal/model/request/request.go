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

// PortalMessageReq 微信推送消息
type PortalMessageReq struct {
	AppID        string `json:"appId" form:"appId"`
	Signature    string `json:"signature" form:"signature"`
	Nonce        string `json:"nonce" form:"nonce"`
	Timestamp    string `json:"timestamp" form:"timestamp"`
	MsgSignature string `json:"msg_signature" form:"msg_signature"`
	EncryptType  string `json:"encryptType" form:"encrypt_type"`
	OpenID       string `json:"openid" form:"openid"`
}

// EncryptMessageReq 加密消息
type EncryptMessageReq struct {
	ToUserName string `json:"ToUserName" form:"ToUserName" xml:"ToUserName"`
	Encrypt    string `json:"encrypt" form:"encrypt" xml:"Encrypt"`
}

type LoginReq struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaKey  string `json:"captchaKey" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}

type PagingQuery struct {
	PageNo   int64 `json:"page_no" form:"page_no"`
	PageSize int64 `json:"page_size" form:"page_size"`
}
