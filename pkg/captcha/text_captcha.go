package captcha

import "github.com/mojocn/base64Captcha"

// IDriver 验证码生成器
type IDriver interface {
	// 画图
	DrawCaptcha(content string) (item base64Captcha.Item, err error)

	// 生成验证码
	GenerateIdQuestionAnswer() (id, q, a string)
}

// IStore 验证码存储
type IStore interface {
	Set(string, string) error

	Get(string, bool) string

	Verify(string, string, bool) bool
}

type TextCaptcha struct {
	*base64Captcha.Captcha
}

// Make 生成验证码
func (t *TextCaptcha) Make() (id, b64s, answer string, err error) {
	return t.Generate()
}

func NewTextCaptcha(store IStore, driver IDriver) *TextCaptcha {
	return &TextCaptcha{Captcha: base64Captcha.NewCaptcha(driver, store)}
}
