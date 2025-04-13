package message

import "github.com/seth16888/wxbusiness/pkg/logger"


func MessageWorker(messageDomain *MessageDomain) <-chan []byte {
	resChan := make(chan []byte)
	go func() {
		defer close(resChan)

		handleMessage(resChan, messageDomain)
	}()
	return resChan
}

// handleMessage 处理消息
//
// 消息类型：
//
//	text 文本消息, image 图片消息, voice 语音消息, video 视频消息, shortvideo 小视频消息,
//	location 地理位置消息, link 链接消息, event 事件消息
func handleMessage(resChan chan<- []byte, domain *MessageDomain) {
	// err := domain.Unmarshal()
	// if err != nil {
	// 	logger.Errorf("message unmarshal error: %s", err.Error())
	// 	return
	// }
	logger.Debugf("-> %s, %s, %s", domain.mixMessage.MsgType, domain.mixMessage.ToUserName,
		domain.mixMessage.FromUserName)

	to := string(domain.mixMessage.FromUserName)
	from := string(domain.mixMessage.ToUserName)
	switch domain.mixMessage.MsgType {
	case MsgTypeText:
		textReply := NewReplyText(to, from, "收到一个文本消息")
		byteSlice, err := domain.MarshalReply(textReply)
		if err != nil {
			logger.Errorf("message marshal error: %s", err.Error())
			return
		}
		resChan <- byteSlice
	case MsgTypeImage:
		textReply := NewReplyText(to, from, "收到一个图片消息")
		byteSlice, err := domain.MarshalReply(textReply)
		if err != nil {
			logger.Errorf("message marshal error: %s", err.Error())
			return
		}
		resChan <- byteSlice
	case MsgTypeVoice:
		textReply := NewReplyText(to, from, "收到一个语音消息")
		byteSlice, err := domain.MarshalReply(textReply)
		if err != nil {
			logger.Errorf("message marshal error: %s", err.Error())
			return
		}
		resChan <- byteSlice
	case MsgTypeVideo:
		textReply := NewReplyText(to, from, "收到一个视频消息")
		byteSlice, err := domain.MarshalReply(textReply)
		if err != nil {
			logger.Errorf("message marshal error: %s", err.Error())
			return
		}
		resChan <- byteSlice
	case MsgTypeShortVideo:
		textReply := NewReplyText(to, from, "收到一个短视频消息")
		byteSlice, err := domain.MarshalReply(textReply)
		if err != nil {
			logger.Errorf("message marshal error: %s", err.Error())
			return
		}
		resChan <- byteSlice
	case MsgTypeLocation:
		textReply := NewReplyText(to, from, "收到一个地理位置消息")
		byteSlice, err := domain.MarshalReply(textReply)
		if err != nil {
			logger.Errorf("message marshal error: %s", err.Error())
			return
		}
		resChan <- byteSlice
	case MsgTypeLink:
		textReply := NewReplyText(to, from, "收到一个链接消息")
		byteSlice, err := domain.MarshalReply(textReply)
		if err != nil {
			logger.Errorf("message marshal error: %s", err.Error())
			return
		}
		resChan <- byteSlice
	case MsgTypeEvent:
		switch domain.mixMessage.Event {
		case EventSubscribe:
			textReply := NewReplyText(to, from, "欢迎关注")
			byteSlice, err := domain.MarshalReply(textReply)
			if err != nil {
				logger.Errorf("message marshal error: %s", err.Error())
				return
			}
			resChan <- byteSlice
		case EventUnsubscribe:
			resChan <- []byte("")
		case EventClick:
			if domain.mixMessage.EventKey != "" {
				textReply := NewReplyText(to, from, "收到一个点击事件消息")
				byteSlice, err := domain.MarshalReply(textReply)
				if err != nil {
					logger.Errorf("message marshal error: %s", err.Error())
					return
				}
				resChan <- byteSlice
			} else {
				textReply := NewReplyText(to, from, "收到了一个Click:"+string(domain.mixMessage.EventKey))
				byteSlice, err := domain.MarshalReply(textReply)
				if err != nil {
					logger.Errorf("message marshal error: %s", err.Error())
					return
				}
				resChan <- byteSlice
			}
		default:
			textReply := NewReplyText(to, from, "收到一个事件消息："+string(domain.mixMessage.Event))
			byteSlice, err := domain.MarshalReply(textReply)
			if err != nil {
				logger.Errorf("message marshal error: %s", err.Error())
				return
			}
			resChan <- byteSlice
		}
	default:
		logger.Debugf("unknown message type: %s", domain.mixMessage.MsgType)
		resChan <- []byte("")
	}
}
