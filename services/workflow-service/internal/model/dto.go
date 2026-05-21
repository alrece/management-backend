package model

// ===== 分类 =====

type CategoryCreateReq struct {
	Name     string `json:"name" binding:"required,max=100"`
	ParentID int64  `json:"parentId"`
	Sort     int    `json:"sort"`
}

type CategoryUpdateReq struct {
	Name     string `json:"name" binding:"required,max=100"`
	ParentID int64  `json:"parentId"`
	Sort     int    `json:"sort"`
}

type CategoryTreeResp struct {
	ID       int64               `json:"id"`
	Name     string              `json:"name"`
	ParentID int64               `json:"parentId"`
	Sort     int                 `json:"sort"`
	Children []*CategoryTreeResp `json:"children"`
}

// ===== 工作流 =====

type WorkflowCreateReq struct {
	Name         string `json:"name" binding:"required,max=200"`
	CategoryID   int64  `json:"categoryId"`
	ParamsSchema string `json:"paramsSchema"`
}

type WorkflowUpdateReq struct {
	Name         string `json:"name" binding:"required,max=200"`
	CategoryID   int64  `json:"categoryId"`
	ParamsSchema string `json:"paramsSchema"`
}

type WorkflowPageReq struct {
	Name       string `form:"name"`
	CategoryID int64  `form:"categoryId"`
	Status     string `form:"status"`
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
}

func (r *WorkflowPageReq) Offset() int {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 10
	}
	return (r.Page - 1) * r.PageSize
}

func (r *WorkflowPageReq) Limit() int {
	if r.PageSize <= 0 {
		return 10
	}
	return r.PageSize
}

type WorkflowResp struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	CategoryID    int64  `json:"categoryId"`
	N8nWorkflowID string `json:"n8nWorkflowId"`
	ParamsSchema  string `json:"paramsSchema"`
	Status        string `json:"status"`
	Creator       int64  `json:"creator"`
	CreatedAt     string `json:"createdAt"`
	Updater       int64  `json:"updater"`
	UpdatedAt     string `json:"updatedAt"`
}

// ===== 执行 =====

type ExecuteReq struct {
	Variables string `json:"variables"` // JSON 字符串
}

type InstancePageReq struct {
	WorkflowID int64  `form:"workflowId"`
	Status     string `form:"status"`
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
}

func (r *InstancePageReq) Offset() int {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 10
	}
	return (r.Page - 1) * r.PageSize
}

func (r *InstancePageReq) Limit() int {
	if r.PageSize <= 0 {
		return 10
	}
	return r.PageSize
}

type InstanceResp struct {
	ID             int64  `json:"id"`
	WorkflowID     int64  `json:"workflowId"`
	N8nExecutionID string `json:"n8nExecutionId"`
	Status         string `json:"status"`
	Variables      string `json:"variables"`
	Result         string `json:"result"`
	DurationMs     int64  `json:"durationMs"`
	ErrorMsg       string `json:"errorMsg"`
	StartedAt      string `json:"startedAt"`
	FinishedAt     string `json:"finishedAt"`
}
