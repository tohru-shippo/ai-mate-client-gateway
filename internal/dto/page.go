package dto

// 用户端分页请求。
type PageQuery struct {
	// 当前页码。
	Current int64 `json:"current"`
	// 每页数量。
	Size int64 `json:"size"`
}

// 用户端分页响应。
type PageResult[T any] struct {
	// 当前页记录。
	Records []T `json:"records"`
	// 当前页码。
	Current int64 `json:"current"`
	// 每页数量。
	Size int64 `json:"size"`
	// 总记录数。
	Total int64 `json:"total"`
	// 总页数。
	Pages int64 `json:"pages"`
}
