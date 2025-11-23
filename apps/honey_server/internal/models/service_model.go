package models

// File: models/service_model.go
// Description: 定义诱捕服务的数据模型及其与镜像、节点的关联关系。

// 服务模型
type ServiceModel struct {
	Model
	Title        string     `gorm:"size:64" json:"title"`        // 服务名称
	Agreement    int8       `json:"agreement"`                   // 协议
	ImageID      uint       `json:"imageID"`                     // 镜像ID
	ImageModel   ImageModel `gorm:"foreignKey:ImageID" json:"-"` // 关联镜像模型
	IP           string     `gorm:"size:32" json:"ip"`           // 服务IP
	Port         int        `json:"port"`                        // 服务端口
	Status       int8       `json:"status"`                      // 服务状态
	HoneyIPCount int        `json:"honeyIPCount"`                // 诱捕IP数量
	ContainerID  string     `gorm:"size:32" json:"containerID"`  // 容器ID
}
