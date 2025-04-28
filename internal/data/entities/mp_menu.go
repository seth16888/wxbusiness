package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

type MPMenu struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"` // MongoDB的主键字段
	AppId           string               `bson:"app_id" json:"app_id"`    // 平台应用ID
	MpId            string               `bson:"mp_id" json:"mp_id"`      // 公众号appid
	MenuID          int64                `bson:"menuid" json:"menuid" `
	Button          []*MenuButton         `bson:"button" json:"button"`
	Conditionalmenu []*ConditionalMenuRes `bson:"conditionalmenu" json:"conditionalmenu"`
	CreatedAt       int64                `bson:"created_at" json:"created_at"`
	UpdatedAt       int64                `bson:"updated_at" json:"updated_at"`
}

// Button 菜单按钮
type MenuButton struct {
	Type       string        `bson:"type" json:"type"`
	Name       string        `bson:"name" json:"name"`
	Key        string        `bson:"key" json:"key"`
	URL        string        `bson:"url" json:"url"`
	MediaID    string        `bson:"media_id" json:"media_id"`
	AppID      string        `bson:"appid" json:"appid"`
	PagePath   string        `bson:"pagepath" json:"pagepath"`
	SubButtons []*MenuButton `bson:"sub_button" json:"sub_button"`
}

// ConditionalMenu 个性化菜单返回结果
type ConditionalMenuRes struct {
	Button    []*MenuButton `bson:"button" json:"button"`
	MatchRule *MatchRule    `bson:"matchrule" json:"matchrule"`
	MenuID    int64        `bson:"menuid" json:"menuid"`
}

// MatchRule 个性化菜单规则
type MatchRule struct {
	TagId              string `bson:"tag_id" json:"tag_id"`
	ClientPlatformType string `bson:"client_platform_type" json:"client_platform_type"`
}
