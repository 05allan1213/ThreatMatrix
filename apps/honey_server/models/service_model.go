package models

// File: models/service_model.go
// Description: 定义诱捕服务的数据模型及其与镜像、节点的关联关系。

import "gorm.io/gorm"

// 服务模型
type ServiceModel struct {
	gorm.Model
	Title        string     `json:"title"`                       // 服务名称
	Agreement    int8       `json:"agreement"`                   // 协议
	ImageID      uint       `json:"imageID"`                     // 镜像ID
	ImageModel   ImageModel `gorm:"foreignKey:ImageID" json:"-"` // 关联镜像模型
	IP           string     `json:"ip"`                          // 服务IP
	Port         int        `json:"port"`                        // 服务端口
	Status       int8       `json:"status"`                      // 服务状态
	HoneyIPCount int        `json:"honeyIPCount"`                // 诱捕IP数量
	ContainerID  string     `json:"containerID"`                 // 容器ID
}
