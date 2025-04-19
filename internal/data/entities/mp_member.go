package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

// MPMember 公众号粉丝
type MPMember struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`                // MongoDB的主键字段
	AppId          string             `bson:"app_id" json:"app_id"`                   // 平台应用ID
	MpId           string             `bson:"mp_id" json:"mp_id"`                     // 公众号appid
	Subscribe      int                `bson:"subscribe" json:"subscribe"`             // 用户是否订阅该公众号标识, 值为0时, 代表此用户没有关注该公众号, 拉取不到其余信息
	OpenId         string             `bson:"openid" json:"openid"`                   // 用户openid
	NickName       string             `bson:"nick_name" json:"nick_name"`             // 用户昵称
	Sex            int                `bson:"sex" json:"sex"`                         // 用户的性别, 值为1时是男性, 值为2时是女性, 值为0时是未知
	Language       string             `bson:"language" json:"language"`               // 用户的语言, 简体中文为zh_CN
	City           string             `bson:"city" json:"city"`                       // 用户所在城市
	Province       string             `bson:"province" json:"province"`               // 用户所在省份
	Country        string             `bson:"country" json:"country"`                 // 用户所在国家
	SubscribeTime  int64              `bson:"subscribe_time" json:"subscribe_time"`   // 用户关注时间, 为时间戳. 如果用户曾多次关注, 则取最后关注时间
	UnionId        string             `bson:"union_id" json:"union_id"`               // 只有在用户将公众号绑定到微信开放平台帐号后, 才会出现该字段.
	Remark         string             `bson:"remark" json:"remark"`                   // 公众号运营者对粉丝的备注, 公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupId        int64              `bson:"group_id" json:"group_id"`               // 用户所在的分组ID（暂时兼容用户分组旧接口）
	Tags           []*MemberTag       `bson:"tags" json:"tags"`                       // 用户被打上的标签ID列表
	SubscribeScene string             `bson:"subscribe_scene" json:"subscribe_scene"` // 用户关注的渠道来源
	QrScene        int64              `bson:"qr_scene" json:"qr_scene"`               // 二维码扫码场景
	QrSceneStr     string             `bson:"qr_scene_str" json:"qr_scene_str"`       // 二维码扫码场景描述
	MessageCount   int64              `bson:"message_count" json:"message_count"`     // 消息发送次数
	CommentCount   int64              `bson:"comment_count" json:"comment_count"`     // 评论次数
	StarComment    int64              `bson:"star_comment" json:"star_comment"`       // 精品留言
	PraiseCount    int64              `bson:"praise_count" json:"praise_count"`       // 点赞数
	PraiseAmounts  int64              `bson:"praise_amounts" json:"praise_amounts"`   // 赞赏总金额：最后两位是小数点后两位，实际金额：10000表示100元
	CreatedAt      int64              `bson:"created_at" json:"created_at"`           // 创建时间
	UpdatedAt      int64              `bson:"updated_at" json:"updated_at"`           // 更新时间
	Blocked        bool               `bson:"blocked" json:"blocked"`                 // 是否被封禁 - 黑名单
}
