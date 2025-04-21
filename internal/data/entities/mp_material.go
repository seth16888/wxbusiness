package entities

import "go.mongodb.org/mongo-driver/bson/primitive"

// MPMaterial 素材
type MPMaterial struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`                      // MongoDB的主键字段
	AppId            string             `bson:"app_id" json:"app_id"`                         // 平台应用ID
	MpId             string             `bson:"mp_id" json:"mp_id"`                           // 公众号appid
	Type             string             `bson:"type" json:"type"`                             // 素材类型: image, voice, video, thumb, ArticleImage, news
	IsPermanent      bool               `bson:"is_permanent" json:"is_permanent"`             // 是否永久素材
	MediaId          string             `bson:"media_id" json:"media_id"`                     // 素材ID
	ThumbMediaId     string             `bson:"thumb_media_id" json:"thumb_media_id"`         // 图文消息的封面图片素材id
	URL              string             `bson:"url" json:"url"`                               // 图文页的URL，或者，当获取的列表是图片素材列表时，该字段是图片的URL
	Filename         string             `bson:"filename" json:"filename"`                     // 文件名
	Path             string             `bson:"path" json:"path"`                             // 本地文件路径
	Title            string             `bson:"title" json:"title"`                           // 图文消息标题
	Author           string             `bson:"author" json:"author"`                         // 图文消息作者
	Digest           string             `bson:"digest" json:"digest"`                         // 图文消息摘要
	ShowCoverPic     int32              `bson:"show_cover_pic" json:"show_cover_pic"`         // 是否显示封面，0为false，即不显示，1为true，即显示
	Content          string             `bson:"content" json:"content"`                       // 图文消息的具体内容，支持HTML标签，必须少于2万字符，小于1M，且此处会去除JS,涉及图片url必须来源 "上传图文消息内的图片获取URL"接口获取。外部图片url将被过滤。
	ContentSourceUrl string             `bson:"content_source_url" json:"content_source_url"` // 图文消息的原文地址，点击“阅读原文”的地址
	CreatedAt        int64              `bson:"created_at" json:"created_at"`
	UpdatedAt        int64              `bson:"updated_at" json:"updated_at"`
}
