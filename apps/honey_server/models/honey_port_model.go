package models

import "gorm.io/gorm"

// 诱捕端口表
type HoneyPortModel struct {
	gorm.Model
	NodeID       uint         `json:"nodeID"`                        // 所属节点ID
	NodeModel    NodeModel    `gorm:"foreignKey:NodeID" json:"-"`    // 关联节点
	NetID        uint         `json:"netID"`                         // 所属网络ID
	NetModel     NetModel     `gorm:"foreignKey:NetID" json:"-"`     // 关联网络
	HoneyIpID    uint         `json:"honeyIpID"`                     // 诱捕IP ID
	HoneyIpModel HoneyIpModel `gorm:"foreignKey:HoneyIpID" json:"-"` // 关联诱捕IP
	ServiceID    uint         `json:"serviceID"`                     // 服务ID
	Port         int          `json:"port"`                          // 服务的端口
	DstIP        string       `gorm:"size:32" json:"dstIP"`          // 目标IP
	DstPort      int          `json:"dstPort"`                       // 目标端口
	Status       int8         `json:"status"`                        // 服务状态
}
