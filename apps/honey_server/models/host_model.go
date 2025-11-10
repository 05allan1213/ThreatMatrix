// Package models 定义诱捕服务所使用的数据实体。
//
// 本文件描述探测到的存活主机模型及其与节点、网络的关系。
package models

import "gorm.io/gorm"

// 存活主机表
type HostModel struct {
	gorm.Model
	NodeID    uint      `json:"nodeID"`                     // 所属节点ID
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"` // 关联节点
	NetID     uint      `json:"netID"`                      // 所属网络ID
	NetModel  NetModel  `gorm:"foreignKey:NetID" json:"-"`  // 关联网络
	IP        string    `gorm:"size:32" json:"ip"`          // 主机IP
	Mac       string    `gorm:"size:64" json:"mac"`         // MAC地址
	Manuf     string    `gorm:"size:64" json:"manuf"`       // 厂商信息
}
