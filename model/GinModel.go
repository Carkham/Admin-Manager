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

type FuncInfo struct {
	NodeName string `json:"node_name"`
	CpuUsage int    `json:"cpu_usage"`
	MemUsage int    `json:"memory_usage"`
	GpuUsage int    `json:"gpu_usage"`
	State    string `json:"state"`
}

type CreateTemplateReq struct {
	TemplateLabel string `json:"template_label"`
	ImageName     string `json:"image_name"`
	BaseCode      string `json:"base_code"`
}

type CreateTemplateResp struct {
	TemplateId int `json:"template_id"`
}
