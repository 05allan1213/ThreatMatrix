package models

// File: models/port_model.go
// Description: 端口模型

type PortModel struct {
	Model
	IP       string `gorm:"size:32" json:"ip"`     // 源ip地址
	Port     int    `json:"port"`                  // 源端口
	DestIp   string `gorm:"size:32" json:"destIp"` // 目标ip地址
	DestPort int    `json:"destPort"`              // 目标端口
}
