package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type PlatformApp struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"` // MongoDB的主键字段
	UserId         string             `bson:"user_id" json:"user_id"`
	Name           string             `bson:"name" json:"name"`
	Type           int64              `bson:"type" json:"type"`
	Token          string             `bson:"token" json:"token"`
	EncodingAesKey string             `bson:"encoding_aes_key" json:"encoding_aes_key"`
	EncodingType   int                `bson:"encoding_type" json:"encoding_type"`
	AppId          string             `bson:"app_id" json:"app_id"` // 公众号appid
	AppSecret      string             `bson:"app_secret" json:"app_secret"`
	Status         int                `bson:"status" json:"status"`
	Introduction   string             `bson:"introduction" json:"introduction"`
	PicUrl         string             `bson:"pic_url" json:"pic_url"`
	CreatedAt      int64              `bson:"created_at" json:"created_at"`
	UpdatedAt      int64              `bson:"updated_at" json:"updated_at"`
	DeletedAt      int64              `bson:"deleted_at" json:"deleted_at"`
	Version        int64              `bson:"version" json:"-"`
}
