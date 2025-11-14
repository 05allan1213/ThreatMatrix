package models

// File: models/matrix_template_model.go
// Description: 定义矩阵模板的数据模型及其与主机模板的关联关系。

// 矩阵模板表
type MatrixTemplateModel struct {
	Model
	Title            string           `gorm:"size:32" json:"title"`                    // 矩阵模板名称
	HostTemplateList HostTemplateList `gorm:"serializer:json" json:"hostTemplateList"` // 主机模板列表
}

type HostTemplateList []HostTemplateInfo

// 主机模板信息
type HostTemplateInfo struct {
	HostTemplateID uint `json:"hostTemplateID"` // 主机模板ID
	Weight         int  `json:"weight"`         // 权重
}
