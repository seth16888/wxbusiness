package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

// MPMaterial 素材
type MPMaterial struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`          // MongoDB的主键字段
	AppId        string             `bson:"app_id" json:"app_id"`             // 平台应用ID
	MpId         string             `bson:"mp_id" json:"mp_id"`               // 公众号appid
	Type         string             `bson:"type" json:"type"`                 // 素材类型
	IsPermanent  bool               `bson:"is_permanent" json:"is_permanent"` // 是否永久素材
  MediaId      string             `bson:"media_id" json:"media_id"`
	ThumbMediaId string             `bson:"thumb_media_id" json:"thumb_media_id"`
	URL          string             `bson:"url" json:"url"`
	Filename     string             `bson:"filename" json:"filename"`
	Path         string             `bson:"path" json:"path"`
	CreatedAt    int64              `bson:"created_at" json:"created_at"`
	UpdatedAt    int64              `bson:"updated_at" json:"updated_at"`
}
