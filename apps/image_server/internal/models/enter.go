package models

// File: models/enter.go
// Description: 基础模型定义，包括通用字段与分页参数结构体

import (
	"time"

	"gorm.io/gorm"
)

// Model 通用字段结构体，包含主键、时间戳和软删除支持
type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`   // 主键 ID
	CreatedAt time.Time      `json:"createdAt"`              // 创建时间
	UpdatedAt time.Time      `json:"updatedAt"`              // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"` // 软删除时间，带索引
}

// PageInfo 通用分页结构体
type PageInfo struct {
	Page  int    `form:"page"`  // 页码
	Limit int    `form:"limit"` // 每页数量
	Key   string `form:"key"`   // 搜索关键字
}

// IDListRequest 批量操作请求结构体
type IDListRequest struct {
	IdList []uint `json:"idList"` // ID 列表
}

// IDRequest 单个 ID 请求结构体
type IDRequest struct {
	ID []uint `json:"id" form:"id" uri:"id"`
}
