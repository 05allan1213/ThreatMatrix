package models

// File: models/node_network_model.go
// Description: 定义节点网卡信息的数据模型及其与节点的关联关系。

import "gorm.io/gorm"

// 节点网卡模型
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
