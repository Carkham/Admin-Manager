package model

type DBFuncInfo struct {
	FunctionID    int64   `gorm:"column:function_id;type:bigint;primaryKey;autoIncrement:true" json:"function_id"` // 函数主键
	FunctionLabel string  `gorm:"column:function_label;type:varchar(128);not null" json:"function_label"`          // 函数名称
	UserID        int64   `gorm:"column:user_id;type:bigint;not null" json:"user_id"`                              // 所属用户的id
	TriggerID     int64   `gorm:"column:trigger_id;type:bigint;not null" json:"trigger_id"`                        // 触发器的id
	SrcType       string  `gorm:"column:src_type;type:varchar(32);not null" json:"src_type"`                       // 源代码导入方式
	SrcLoc        string  `gorm:"column:src_loc;type:varchar(256);not null" json:"src_loc"`                        // 代码定位数据
	Replicas      int32   `gorm:"column:replicas;type:int;not null;default:1" json:"replicas"`                     // 副本数
	QuotaInfo     string  `gorm:"column:quota_info;type:text;not null" json:"quota_info"`                          // 配额信息
	TemplateID    int64   `gorm:"column:template_id;type:bigint;not null" json:"template_id"`                      // 模板ID
	ImageName     string  `gorm:"column:image_name;type:varchar(256);not null" json:"image_name"`                  // 模版对应镜像定位符号
	TemplateLabel string  `gorm:"column:template_label;type:varchar(128);not null" json:"template_label"`          // 模板标签
	FileName      string  `gorm:"column:file_name;type:varchar(256);not null" json:"file_name"`                    // 工程模板文件名
	TriggerType   string  `gorm:"column:trigger_type;type:varchar(64);not null" json:"trigger_type"`               // 触发器类型
	TriggerConfig string  `gorm:"column:trigger_config;type:text;not null" json:"trigger_config"`                  // 触发器配置信息
	TriggerLabel  *string `gorm:"column:trigger_label;type:varchar(256)" json:"trigger_label"`                     // 触发器标签
}

type QuotaInfo struct {
	CpuRequest int `json:"cpu_request"`
	MemRequest int `json:"mem_request"`
	CpuLimit   int `json:"cpu_limit"`
	MemLimit   int `json:"mem_limit"`
	GpuQuota   int `json:"gpu_quota"`
}

type CronJobConfig struct {
	TrigWeekday int `json:"trig_weekday"`
	TrigHour    int `json:"trig_hour"`
	TrigMin     int `json:"trig_min"`
}

type DBFuncOverview struct {
	FunctionID    int64  `gorm:"column:function_id;type:bigint;primaryKey;autoIncrement:true" json:"function_id"` // 函数主键
	FunctionLabel string `gorm:"column:function_label;type:varchar(128);not null" json:"function_label"`          // 函数名称
	TemplateLabel string `gorm:"column:template_label;type:varchar(128);not null" json:"template_label"`          // 模板标签
	TriggerType   string `gorm:"column:trigger_type;type:varchar(64);not null" json:"trigger_type"`               // 触发器类型
}
