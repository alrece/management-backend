package page

// Req 分页请求基类
type Req struct {
	Page     int `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" json:"pageSize" binding:"omitempty,min=1,max=100"`
}

func (r *Req) GetPage() int {
	if r.Page <= 0 {
		return 1
	}
	return r.Page
}

func (r *Req) GetPageSize() int {
	if r.PageSize <= 0 {
		return 10
	}
	return r.PageSize
}

func (r *Req) Offset() int {
	return (r.GetPage() - 1) * r.GetPageSize()
}

func (r *Req) Limit() int {
	return r.GetPageSize()
}
