package model

type JobTaskCreateReq struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Handler  string `json:"handler" binding:"required,min=1,max=100"`
	CronExpr string `json:"cronExpr" binding:"required"`
	Params   string `json:"params"`
	Status   *int8  `json:"status"`
	Remark   string `json:"remark"`
}

type JobTaskUpdateReq struct {
	Name     string `json:"name" binding:"omitempty,min=1,max=100"`
	Handler  string `json:"handler" binding:"omitempty,min=1,max=100"`
	CronExpr string `json:"cronExpr" binding:"omitempty"`
	Params   string `json:"params"`
	Status   *int8  `json:"status"`
	Remark   string `json:"remark"`
}

type JobTaskPageReq struct {
	Name   string `form:"name"`
	Status *int8  `form:"status"`
	Page   int    `form:"page"`
	Size   int    `form:"pageSize"`
}

func (r *JobTaskPageReq) Offset() int {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	return (r.Page - 1) * r.Size
}

func (r *JobTaskPageReq) Limit() int {
	if r.Size <= 0 {
		return 10
	}
	return r.Size
}

type JobTaskResp struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Handler  string `json:"handler"`
	CronExpr string `json:"cronExpr"`
	Params   string `json:"params"`
	Status   int8   `json:"status"`
	Remark   string `json:"remark"`
	Creator  *int64 `json:"creator"`
}

type ExecLogPageReq struct {
	TaskID     *int64 `form:"taskId"`
	Status     *int8  `form:"status"`
	StartTime  string `form:"startTime"`
	EndTime    string `form:"endTime"`
	Page       int    `form:"page"`
	Size       int    `form:"pageSize"`
}

func (r *ExecLogPageReq) Offset() int {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Size <= 0 {
		r.Size = 10
	}
	return (r.Page - 1) * r.Size
}

func (r *ExecLogPageReq) Limit() int {
	if r.Size <= 0 {
		return 10
	}
	return r.Size
}

type ExecLogResp struct {
	ID          int64   `json:"id"`
	TaskID      int64   `json:"taskId"`
	TaskName    string  `json:"taskName"`
	TriggerType int8    `json:"triggerType"`
	Status      int8    `json:"status"`
	DurationMs  int64   `json:"durationMs"`
	Result      string  `json:"result"`
	Error       string  `json:"error"`
	StartTime   string  `json:"startTime"`
}
