package request

// CreateQRCodeReq 创建公众号二维码
type CreateQRCodeReq struct {
	Exp   int    `json:"expire_seconds" binding:"omitempty,min=60,max=2592000" msg:"expire_seconds,min=60,max=2592000"`
	Scene string `json:"scene" binding:"required" msg:"scene required"`
}
