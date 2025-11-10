// Package models 定义诱捕服务所使用的数据实体。
//
// 本文件描述主机模板及其端口配置的序列化结构。
package models

import "gorm.io/gorm"

// 主机模板表
type HostTemplateModel struct {
	gorm.Model
	Title    string               `gorm:"size:32" json:"title"`            // 主机名称
	PortList HostTemplatePortList `gorm:"serializer:json" json:"portList"` // 主机端口列表
}

type HostTemplatePortList []HostTemplatePort

// 主机模板端口列表
type HostTemplatePort struct {
	Port      int  `json:"port"`      // 端口号
	ServiceID uint `json:"serviceID"` // 关联服务ID
}
