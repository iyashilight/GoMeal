package utils

// PageReq 分页请求参数，支持从 Query String 绑定
type PageReq struct {
	Page int `json:"page" form:"page"` // 页码（从1开始）
	Size int `json:"size" form:"size"` // 每页条数
}

// Normalize 规范化分页参数，设置默认值和上限
func (p *PageReq) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	if p.Size > 50 {
		p.Size = 50
	}
}

// PageResp 分页响应结构
type PageResp struct {
	Total int64       `json:"total"` // 总记录数
	Page  int         `json:"page"`  // 当前页码
	Size  int         `json:"size"`  // 每页条数
	Items interface{} `json:"items"` // 当前页数据
}
