package model

import "time"

type ResourceQuota struct {
	DownCPU int
	UpCPU   int
	DownMem int
	UpMem   int
	GPU     int
}

type PodInfo struct {
	Name string
	Node string
}

type DeploymentInfo struct {
	Name      string
	FuncIDStr string
	Replicas  int
	Status    string
}

type PodMetricInfo struct {
	FunctionID int64
	DeployName string
	PodName    string
	NodeName   string
	State      string
	CpuUsage   int
	MemUsage   int
	GpuUsage   int
}

type NodeInfo struct {
	Name     string
	Status   string
	Optional bool
	Age      string
	Version  string
	CpuTotal int
	MemTotal int
	GpuTotal int
	CpuUse   int
	MemUse   int
	GpuUse   int
}

type FuncInfo struct {
	FunctionID     int64
	UserID         int64
	UserName       string
	FuncLabel      string
	SourceType     string
	SourceLocation string
	TrigType       string
	TimeStr        string
	PodCount       int32
	ImageName      string
	CPUQuotaM      [2]int
	MemQuotaMi     [2]int
	GPUQuota       int
}

type PodMetricResp struct {
	Kind       string          `json:"kind"`
	APIVersion string          `json:"apiVersion"`
	Items      []PodMetricItem `json:"items"`
}

type PodMetricItem struct {
	Metadata   PodMetadata     `json:"metadata"`
	Timestamp  time.Time       `json:"timestamp"`
	Window     string          `json:"window"`
	Containers []PodContainers `json:"containers"`
}

type PodMetadata struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
}

type PodUsage struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type PodContainers struct {
	Name  string   `json:"name"`
	Usage PodUsage `json:"usage"`
}

type NodeMetricResp struct {
	Kind       string            `json:"kind"`
	APIVersion string            `json:"apiVersion"`
	Items      []NodeMetricItems `json:"items"`
}

type NodeMetadata struct {
	Name              string            `json:"name"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
}
type NodeUsage struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
type NodeMetricItems struct {
	Metadata  NodeMetadata `json:"metadata"`
	Timestamp time.Time    `json:"timestamp"`
	Window    string       `json:"window"`
	Usage     NodeUsage    `json:"usage"`
}
