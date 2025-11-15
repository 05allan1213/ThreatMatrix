package models

// File: models/honey_port_model.go
// Description: 蜜罐端口模型

// 蜜罐端口模型
type HoneyPortModel struct {
	Model
	NodeID    uint   `json:"nodeID"`               // 节点id
	NetID     uint   `json:"netID"`                // 网络id
	HoneyIpID uint   `json:"honeyIpID"`            // 蜜罐ipid
	ServiceID uint   `json:"serviceID"`            // 服务id
	Port      int    `json:"port"`                 // 服务的端口
	DstIP     string `gorm:"size:32" json:"dstIP"` // 目标ip
	DstPort   int    `json:"dstPort"`              // 目标端口
	Status    int8   `json:"status"`
}
