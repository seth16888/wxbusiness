package message

import "encoding/xml"

type MixMessage struct {
	CommonToken

	// 基本消息
	MsgId         int64   `xml:"MsgId"` // 其他消息推送过来是 MsgId
	TemplateMsgId int64   `xml:"MsgID"` // 模板消息推送成功的消息是 MsgID
	Content       string  `xml:"Content"`
	Recognition   string  `xml:"Recognition"`
	PicURL        string  `xml:"PicUrl"`
	MediaId       string  `xml:"MediaId"`
	Format        string  `xml:"Format"`
	ThumbMediaId  string  `xml:"ThumbMediaId"`
	LocationX     float64 `xml:"Location_X"`
	LocationY     float64 `xml:"Location_Y"`
	Scale         float64 `xml:"Scale"`
	Label         string  `xml:"Label"`
	Title         string  `xml:"Title"`       // 链接消息：标题
	Description   string  `xml:"Description"` // 链接消息：描述
	URL           string  `xml:"Url"`         // 链接消息：链接
	BizMsgMenuId  int64   `xml:"bizmsgmenuid"`
	MsgDataId     string  `xml:"MsgDataId"`  // 消息来自来自文章
	Idx           int     `xml:"Idx"`        // 多图文时第几篇文章，从1开始（消息如果来自文章时才有）
	MediaId16K    string  `xml:"MediaId16K"` // 16K采样率语音消息媒体id，可以调用获取临时素材接口拉取数据

	// 事件相关
	Event       EventType `xml:"Event" json:"Event"`
	EventKey    string    `xml:"EventKey"`  // 事件KEY值，qrscene_为前缀，后面为二维码的场景值ID
	Ticket      string    `xml:"Ticket"`    // 二维码的ticket，可用来换取二维码图片
	Latitude    string    `xml:"Latitude"`  // 上报地理位置：纬度
	Longitude   string    `xml:"Longitude"` // 上报地理位置：经度
	Precision   string    `xml:"Precision"` // 上报地理位置：精度
	MenuId      string    `xml:"MenuId"`
	Status      string    `xml:"Status"`
	SessionFrom string    `xml:"SessionFrom"`
	TotalCount  int64     `xml:"TotalCount"`
	FilterCount int64     `xml:"FilterCount"`
	SentCount   int64     `xml:"SentCount"`
	ErrorCount  int64     `xml:"ErrorCount"`

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"`
		ScanResult string `xml:"ScanResult"`
	} `xml:"ScanCodeInfo"`

	SendPicsInfo struct {
		Count   int32      `xml:"Count"`
		PicList []EventPic `xml:"PicList>item"`
	} `xml:"SendPicsInfo"`

	SendLocationInfo struct {
		LocationX float64 `xml:"Location_X"`
		LocationY float64 `xml:"Location_Y"`
		Scale     float64 `xml:"Scale"`
		Label     string  `xml:"Label"`
		Poiname   string  `xml:"Poiname"`
	}

	subscribeMsgPopupEventList []SubscribeMsgPopupEvent `json:"-"`

	// 订阅通知消息
	SubscribeMsgPopupEvent []struct {
		List SubscribeMsgPopupEvent `xml:"List"`
	} `xml:"SubscribeMsgPopupEvent"`
}

// SetSubscribeMsgPopupEvents 设置订阅消息事件
func (s *MixMessage) SetSubscribeMsgPopupEvents(list []SubscribeMsgPopupEvent) {
	s.subscribeMsgPopupEventList = list
}

// GetSubscribeMsgPopupEvents 获取订阅消息事件数据
func (s *MixMessage) GetSubscribeMsgPopupEvents() []SubscribeMsgPopupEvent {
	if s.subscribeMsgPopupEventList != nil {
		return s.subscribeMsgPopupEventList
	}
	list := make([]SubscribeMsgPopupEvent, len(s.SubscribeMsgPopupEvent))
	for i, item := range s.SubscribeMsgPopupEvent {
		list[i] = item.List
	}
	return list
}

// CommonToken 消息中通用的结构
type CommonToken struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATA    `xml:"ToUserName" json:"ToUserName"`
	FromUserName CDATA    `xml:"FromUserName" json:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime" json:"CreateTime"`
	MsgType      MsgType  `xml:"MsgType" json:"MsgType"`
}

// SetToUserName set ToUserName
func (msg *CommonToken) SetToUserName(toUserName CDATA) {
	msg.ToUserName = toUserName
}

// SetFromUserName set FromUserName
func (msg *CommonToken) SetFromUserName(fromUserName CDATA) {
	msg.FromUserName = fromUserName
}

// SetCreateTime set createTime
func (msg *CommonToken) SetCreateTime(createTime int64) {
	msg.CreateTime = createTime
}

// SetMsgType set MsgType
func (msg *CommonToken) SetMsgType(msgType MsgType) {
	msg.MsgType = msgType
}

// GetOpenID get the FromUserName value
func (msg *CommonToken) GetOpenID() string {
	return string(msg.FromUserName)
}
