package models

// File: models/honey_ip_model.go
// Description: 定义诱捕 IP 的数据模型及其与节点、网络的关联关系。

// 诱捕IP表
type HoneyIpModel struct {
	Model
	NodeID    uint      `json:"nodeID"`                         // 所属节点ID
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"`     // 关联节点
	NetID     uint      `gorm:"index:idx_net_id" json:"netID"`  // 所属网络ID
	NetModel  NetModel  `gorm:"foreignKey:NetID" json:"-"`      // 关联网络
	IP        string    `gorm:"size:32;index:idx_ip" json:"ip"` // 诱捕IP
	Mac       string    `gorm:"size:64" json:"mac"`             // MAC地址
	Network   string    `gorm:"size:32" json:"network"`         // 网卡
	Status    int8      `json:"status"`                         // 状态  1 创建中 2 运行中 3 失败 4 删除中
	ErrorMsg  string    `gorm:"size:64" json:"errorMsg"`        // 错误信息
}
