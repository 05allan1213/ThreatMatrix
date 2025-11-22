package models

// File: models/ip_model.go
// Description: IP网络接口模型

// IpModel IP网络接口模型
type IpModel struct {
	Model           // 嵌入基础模型
	Ip       string `json:"ip"`       // IP地址
	Mask     int8   `json:"mask"`     // 子网掩码位数
	LinkName string `json:"linkName"` // 节点上创建的macvlan子接口名称
	Network  string `json:"network"`  // 网卡名称
	Mac      string `json:"mac"`      // 网卡的MAC地址
}
