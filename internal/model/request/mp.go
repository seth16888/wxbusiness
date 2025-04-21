package request

// CreateQRCodeReq 创建公众号二维码
type CreateQRCodeReq struct {
	Exp   int    `json:"expire_seconds" binding:"omitempty,min=60,max=2592000" msg:"expire_seconds,min=60,max=2592000"`
	Scene string `json:"scene" binding:"required" msg:"scene required"`
}

// PullMaterialReq 拉取永久素材
type PullMaterialReq struct {
	Type   string `json:"type" binding:"required" msg:"type required"`
	Offset int64  `json:"offset" binding:"required" msg:"offset required"`
	Count  int64  `json:"count" binding:"required" msg:"count required"`
}
