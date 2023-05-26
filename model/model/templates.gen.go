// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameTemplate = "templates"

// Template mapped from table <templates>
type Template struct {
	TemplateID    int64  `gorm:"column:template_id;type:bigint;primaryKey;autoIncrement:true" json:"template_id"` // 主键id
	ImageName     string `gorm:"column:image_name;type:varchar(256);not null" json:"image_name"`                  // 模版对应镜像定位符号
	TemplateLabel string `gorm:"column:template_label;type:varchar(128);not null" json:"template_label"`          // 模板标签
	FileName      string `gorm:"column:file_name;type:varchar(256);not null" json:"file_name"`                    // 工程模板文件名
}

// TableName Template's table name
func (*Template) TableName() string {
	return TableNameTemplate
}
