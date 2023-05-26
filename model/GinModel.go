package model

type JSONResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Extra   interface{} `json:"extra"`
}

type GetListResp struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

type NodeItem struct {
	Name     string `json:"node_name"`
	Status   string `json:"status"`
	Optional bool   `json:"optional"`
	Age      string `json:"age"`
	Version  string `json:"version"`
	CpuUsage string `json:"cpu_usage"`
	CpuTotal string `json:"cpu_total"`
	MemUsage string `json:"memory_usage"`
	MemTotal string `json:"memory_total"`
	GpuUsage string `json:"gpu_usage"`
	GpuTotal string `json:"gpu_total"`
}

type FuncList struct {
	UserName     string      `json:"user_name"`
	FunctionId   int         `json:"function_id"`
	FunctionName string      `json:"function_name"`
	TemplateName string      `json:"template_name"`
	State        string      `json:"state"`
	ReplicasInfo interface{} `json:"replicas_info"`
}

type CreateTemplateReq struct {
	TemplateLabel string `json:"template_label"`
	ImageName     string `json:"image_name"`
	BaseCode      string `json:"base_code"`
}

type CreateTemplateResp struct {
	TemplateId int `json:"template_id"`
}

type UserInfo struct {
	UserId    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	LastLogin string `json:"last_login"`
	IsAdmin   bool   `json:"is_admin"`
	Email     string `json:"email"`
}

type CreateUserReq struct {
	Username string `json:"user_name"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
	Password string `json:"password"`
}

type CreateUserResp struct {
	UserId int64 `json:"user_id"`
}
