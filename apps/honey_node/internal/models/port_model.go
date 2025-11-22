package models

// File: models/port_model.go
// Description: 端口模型

type PortModel struct {
	Model
	LocalAddr  string `gorm:"size:64" json:"localAddr"`  // 本地地址
	TargetAddr string `gorm:"size:64" json:"targetAddr"` // 目标地址
}
