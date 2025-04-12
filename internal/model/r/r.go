package r

import "strconv"

// R Server 返回给客户端的数据结构
type R[T any] struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Data    T      `json:"data,omitempty"`
}

// NewR 返回一个新的 R 实例
func NewR[T any](code int64, message string, data T) *R[T] {
	c := "00000"
	if code != 0 {
		c = strconv.FormatInt(code, 10)
	}
	return &R[T]{
		Code:    c,
		Message: message,
		Data:    data,
	}
}

// Success 返回成功的 R 实例
func Success[T any](data T) *R[T] {
	return NewR(0, "success", data)
}

// Error 返回错误的 R 实例
func Error[T any](code int64, message string) *R[T] {
	var zero T
	return NewR(code, message, zero)
}
