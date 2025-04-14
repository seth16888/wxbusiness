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
