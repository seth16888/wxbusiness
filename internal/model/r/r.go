package r

// R Server 返回给客户端的数据结构
type R struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`
	Data    any    `json:"data,omitempty"`
}

func (r *R) IsSuccess() bool {
	return r.Code == 0
}

func (r *R) StatusCode() int {
  if r.Code == 0 {
    return 200
  }
  // 1000 以下的 code 直接返回
  // 比如：400-409 500-511
  if r.Code < 1000 {
    return int(r.Code)
  }
  // 1000 以上的 是业务错误
  // 直接返回 200
  return 200
}

// NewR 返回一个新的 R 实例
func NewR(code int64, message string, data any) *R {
	return &R{Code: code, Message: message, Data: data}
}

// Success 返回成功的 R 实例
func Success() *R {
	return NewR(200, "success", nil)
}

// SuccessData 返回成功的 R 实例，带数据
func SuccessData(data any) *R {
	return NewR(200, "success", data)
}

// Error 返回错误的 R 实例
func Error(code int64, message string) *R {
	return NewR(code, message, nil)
}

