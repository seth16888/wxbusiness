package response

type CreateAppResp struct {
	AppId string `json:"appId"`
}

// GetCaptchaResp 获取验证码
type GetCaptchaResp struct {
	CaptchaKey   string `json:"captchaKey"`
	CaptchaValue string `json:"captchaValue"`
}

// LoginResp 登录响应
type LoginResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64    `json:"expiresIn"`
}

// Ticket 获取二维码返回结果
type Ticket struct {
	Ticket        string `json:"ticket"`
	URL           string `json:"url"` // URL 解析后的网址，可根据URL自行生成二维码
	ExpireSeconds int64  `json:"expire_seconds"`
}
