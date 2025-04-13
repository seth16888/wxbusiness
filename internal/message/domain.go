package message

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

type MessageDomain struct {
	// app config
	// app *entities.PlatformApp

	// encryptType
	encryptType string
	// contentType
	contentType string

	// OpenId
	OpenId string
	// RandomString
	RandomString string
	// Timestamp
	Timestamp string
	// Nonce
	Nonce string

	// IsXML
	IsXML bool
	// IsEncrypted
	IsEncrypted bool

	// DataReader
	dataReader io.ReadCloser

	// MixMessage
	mixMessage *MixMessage
}

func NewMessageDomain(body io.ReadCloser) *MessageDomain {
	return &MessageDomain{
		// app:         appConf,
		dataReader:  body,
		IsXML:       false,
		IsEncrypted: false,
		contentType: "text/xml",
		mixMessage:  new(MixMessage),
	}
}

// SetEncryptType
func (domain *MessageDomain) SetEncryptType(encryptType string) {
	domain.encryptType = encryptType
	domain.IsEncrypted = encryptType != ""
}

// SetContentType
func (domain *MessageDomain) SetContentType(contentType string) {
	domain.contentType = contentType
	domain.IsXML = contentType == "text/xml" || contentType == "application/xml"
}

// Unmarshal
func (domain *MessageDomain) Unmarshal() error {
	rawDataBytes, err := io.ReadAll(domain.dataReader)
	if err != nil {
		return fmt.Errorf("从body中读取数据失败, err=%v", err)
	}

	msg := new(MixMessage)
	if domain.IsXML {
		if err := xml.Unmarshal(rawDataBytes, msg); err != nil {
			return fmt.Errorf("从body中解析xml失败,err=%v", err)
		}
	} else {
		if err := json.Unmarshal(rawDataBytes, msg); err != nil {
			return fmt.Errorf("从body中解析json失败,err=%v", err)
		}
	}

	domain.mixMessage = msg

	return nil
}

// MarshalReply
func (domain *MessageDomain) MarshalReply(reply interface{}) ([]byte, error) {
	bytes, err := xml.Marshal(reply)
	if err != nil {
		return []byte(""), fmt.Errorf("回复消息序列化失败, err=%v", err)
	}

	return bytes, nil
}
