package model

type PageResult[T any] struct {
	List  []T   `json:"list"`
	Total int64 `json:"total"`
}

func GetLimitOffset(pageNum int, pageSize int) (offset int, limit int) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset = (pageNum - 1) * pageSize
	limit = pageSize
	return
}

// NewPageResult 构造函数，用于创建一个新的 PageResult 实例。
func NewPageResult[T any]() *PageResult[T] {
	return &PageResult[T]{
		List:  make([]T, 0),
		Total: 0,
	}
}
