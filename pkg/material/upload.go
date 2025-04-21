package material

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/seth16888/wxcommon/hc"
)

// MultipartFormField 保存文件或其他字段信息
type MultipartFormField struct {
	IsFile     bool
	Fieldname  string
	Value      []byte
	FilePath   string
	Filename   string
	FileReader io.Reader
}

func UploadFile(fieldName, filePath, uri string, hc *hc.Client) ([]byte, error) {
	return PostFileByFormData(uri, nil, filePath, fieldName, hc)
}

// 使用multipart/form-data发送文件
func PostFileByFormData(url string, params map[string]string, file_path string, file_key string, hc *hc.Client) (respBody []byte, err error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, errors.New("PostFileByFormData open file error")
	}
	defer file.Close()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	for key, val := range params {
		_ = bodyWriter.WriteField(key, val)
	}

	fileWrite, err := bodyWriter.CreateFormFile(file_key, file.Name())
	if err != nil {
		return nil, errors.New("PostFileByFormData create part error")
	}
	_, err = io.Copy(fileWrite, file)
	if err != nil {
		return nil, errors.New("PostFileByFormData read file error")
	}

	if err = bodyWriter.Close(); err != nil {
		return nil, errors.New("PostFileByFormData close multipart error")
	}

	contentType := bodyWriter.FormDataContentType()
	dataBytesReader := bytes.NewReader(bodyBuf.Bytes())

	resp, e := hc.Post(url, contentType, dataBytesReader)
	if e != nil {
		err = fmt.Errorf("upload media to wx error,%s", e.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error : url=%s , statusCode=%d", url, resp.StatusCode)
	}

	respBody, err = io.ReadAll(resp.Body)
	return
}

// PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(fields []MultipartFormField, uri string, hc *hc.Client) (respBody []byte, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%v", e)
				return nil, err
			}

			if field.FileReader == nil {
				fh, e := os.Open(field.FilePath)
				if e != nil {
					err = fmt.Errorf("error opening file , err=%v", e)
					return nil, err
				}
				_, err = io.Copy(fileWriter, fh)
				_ = fh.Close()
				if err != nil {
					return nil, err
				}
			} else {
				if _, err = io.Copy(fileWriter, field.FileReader); err != nil {
					return nil, err
				}
			}
		} else {
			partWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			if e != nil {
				err = e
				return nil, err
			}
			valueReader := bytes.NewReader(field.Value)
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return nil, err
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, e := hc.Post(uri, contentType, bodyBuf)
	if e != nil {
		err = e
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error : uri=%v , statusCode=%v", uri, resp.StatusCode)
	}
	respBody, err = io.ReadAll(resp.Body)
	return respBody, err
}
