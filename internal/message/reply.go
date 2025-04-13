package message

import "time"

// ReplyImage 图片回复
type ReplyImage struct {
	CommonToken

	Image struct {
		MediaId string `xml:"MediaId"`
	} `xml:"Image"`
}

func NewReplyImage(to, from, mediaId string) *ReplyImage {
	reply := new(ReplyImage)
	reply.SetToUserName(CDATA(to))
	reply.SetFromUserName(CDATA(from))
	reply.SetCreateTime(time.Now().Unix())
	reply.SetMsgType(MsgTypeImage)
	reply.Image.MediaId = mediaId
	return reply
}

// ReplyText 文本回复
type ReplyText struct {
	CommonToken
	Content CDATA `xml:"Content"`
}

func NewReplyText(to, from, content string) *ReplyText {
	reply := new(ReplyText)
	reply.SetToUserName(CDATA(to))
	reply.SetFromUserName(CDATA(from))
	reply.SetMsgType(MsgTypeText)
	reply.SetCreateTime(time.Now().Unix())
	reply.Content = CDATA(content)
	return reply
}

// ReplyVoice 语音回复
type ReplyVoice struct {
	CommonToken
	Voice struct {
		MediaId CDATA `xml:"MediaId"`
	} `xml:"Voice"`
}

func NewReplyVoice(to, from, mediaId string) *ReplyVoice {
	reply := new(ReplyVoice)
	reply.SetToUserName(CDATA(to))
	reply.SetFromUserName(CDATA(from))
	reply.SetMsgType(MsgTypeVoice)
	reply.SetCreateTime(time.Now().Unix())
	reply.Voice.MediaId = CDATA(mediaId)
	return reply
}

// ReplyVideo 视频回复
type ReplyVideo struct {
	CommonToken
	Video struct {
		MediaId     CDATA `xml:"MediaId"`
		Title       CDATA `xml:"Title"`
		Description CDATA `xml:"Description"`
	} `xml:"Video"`
}

func NewReplyVideo(to, from, mediaId, title, description string) *ReplyVideo {
	reply := new(ReplyVideo)
	reply.SetToUserName(CDATA(to))
	reply.SetFromUserName(CDATA(from))
	reply.SetMsgType(MsgTypeVideo)
	reply.SetCreateTime(time.Now().Unix())
	reply.Video.MediaId = CDATA(mediaId)
	reply.Video.Title = CDATA(title)
	reply.Video.Description = CDATA(description)
	return reply
}

// ReplyMusic 音乐回复
type ReplyMusic struct {
	CommonToken
	Music struct {
		Title        CDATA `xml:"Title"`
		Description  CDATA `xml:"Description"`
		MusicURL     CDATA `xml:"MusicURL"`
		HQMusicURL   CDATA `xml:"HQMusicUrl"`
		ThumbMediaId CDATA `xml:"ThumbMediaId"`
	} `xml:"Music"`
}

func NewReplyMusic(to, from, title, description, musicURL, hqMusicURL, thumbMediaId string) *ReplyMusic {
	reply := new(ReplyMusic)
	reply.SetToUserName(CDATA(to))
	reply.SetFromUserName(CDATA(from))
	reply.SetMsgType(MsgTypeMusic)
	reply.SetCreateTime(time.Now().Unix())
	reply.Music.Title = CDATA(title)
	reply.Music.Description = CDATA(description)
	reply.Music.MusicURL = CDATA(musicURL)
	reply.Music.HQMusicURL = CDATA(hqMusicURL)
	reply.Music.ThumbMediaId = CDATA(thumbMediaId)
	return reply
}

// ReplyNewsItem 回复图文消息的子项
type ReplyNewsItem struct {
	Title       CDATA `xml:"Title"`
	Description CDATA `xml:"Description"`
	PicUrl      CDATA `xml:"PicUrl"`
	Url         CDATA `xml:"Url"`
}

// ReplyNews 回复图文消息
type ReplyNews struct {
	CommonToken
	ArticleCount int `xml:"ArticleCount"`
	Articles     struct {
		Item []ReplyNewsItem `xml:"item"`
	} `xml:"Articles"`
}

func NewReplyNews(to, from string, articles []ReplyNewsItem) *ReplyNews {
	reply := new(ReplyNews)
	reply.SetToUserName(CDATA(to))
	reply.SetFromUserName(CDATA(from))
	reply.SetCreateTime(time.Now().Unix())
	reply.SetMsgType(MsgTypeNews)
	reply.ArticleCount = len(articles)
	reply.Articles.Item = articles
	return reply
}
