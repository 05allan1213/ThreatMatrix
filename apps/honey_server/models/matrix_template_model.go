// Package models 定义诱捕服务所使用的数据实体。
//
// 本文件描述矩阵模板及其包含的主机模板映射关系。
package models

import "gorm.io/gorm"

// 矩阵模板表
type MatrixTemplateModel struct {
	gorm.Model
	Title            string           `gorm:"size:32" json:"title"`                    // 矩阵模板名称
	HostTemplateList HostTemplateList `gorm:"serializer:json" json:"hostTemplateList"` // 主机模板列表
}

type HostTemplateList []HostTemplateInfo

// 主机模板信息
type HostTemplateInfo struct {
	HostTemplateID uint `json:"hostTemplateID"` // 主机模板ID
	Weight         int  `json:"weight"`         // 权重
}
