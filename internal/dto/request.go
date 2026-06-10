package dto

// 用户端通用查询/更新请求结构。
type Request struct {
	// 查询条件。
	Filter map[string]any `json:"filter"`
	// 更新字段。
	Update map[string]any `json:"update"`
}

// 用户端通用分页请求结构。
type PageRequest struct {
	// 查询条件。
	Filter map[string]any `json:"filter"`
	// 分页参数。
	Page PageQuery `json:"page"`
}
