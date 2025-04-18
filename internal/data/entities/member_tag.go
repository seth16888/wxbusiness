package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

// MemberTag 会员标签
// MongoDB数据库表名：member_tags
type MemberTag struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"` // MongoDB的主键字段
	AppId     string             `bson:"app_id" json:"app_id"`    // 平台应用ID
	MpId      string             `bson:"mp_id" json:"mp_id"`      // 公众号appid
	TagId     int64              `bson:"tag_id" json:"tag_id"`    // 微信标签ID
	Name      string             `bson:"name" json:"name"`
	Count     int64              `bson:"count" json:"count"` // 标签下粉丝数
	CreatedAt int64              `bson:"created_at" json:"created_at"`
	UpdatedAt int64              `bson:"updated_at" json:"updated_at"`
}
