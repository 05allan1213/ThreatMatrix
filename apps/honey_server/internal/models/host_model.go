package models

// File: models/host_model.go
// Description: 定义探测到的存活主机模型及其与节点、网络的关系。

// 存活主机表
type HostModel struct {
	Model
	NodeID    uint      `json:"nodeID"`                         // 所属节点ID
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"`     // 关联节点
	NetID     uint      `gorm:"index:idx_net_id" json:"netID"`  // 所属网络ID
	NetModel  NetModel  `gorm:"foreignKey:NetID" json:"-"`      // 关联网络
	IP        string    `gorm:"size:32;index:idx_ip" json:"ip"` // 主机IP
	Mac       string    `gorm:"size:64" json:"mac"`             // MAC地址
	Manuf     string    `gorm:"size:64" json:"manuf"`           // 厂商信息
}
