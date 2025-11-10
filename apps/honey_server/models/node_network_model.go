// Package models 定义诱捕服务所使用的数据实体。
//
// 本文件描述节点网卡信息的模型，用于记录各接口的网络属性。
package models

import "gorm.io/gorm"

// 节点网卡表
type NodeNetworkModel struct {
	gorm.Model
	NodeID    uint      `json:"nodeID"`                     // 关联的节点ID
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"` // 关联的节点模型
	Network   string    `gorm:"size:32" json:"network"`     // 网卡名称
	IP        string    `gorm:"size:32" json:"ip"`          // 探针IP
	Mask      string    `gorm:"size:32" json:"mask"`        // 子网掩码 8-32
	Gateway   string    `gorm:"size:32" json:"gateway"`     // 网关
	Status    int8      `json:"status"`                     // 网关状态
}
