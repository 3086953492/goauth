package dto

import "io"

type PaginationResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`      // 总记录数
	Page       int   `json:"page"`       // 当前页
	PageSize   int   `json:"pageSize"`   // 每页大小
	TotalPages int   `json:"totalPages"` // 总页数
}

type FileMeta struct {
	Data        io.Reader
	Filename    string
	ContentType string
	Size        int64
}
