// Package models 定义诱捕服务所使用的数据实体。
//
// 本文件描述诱捕服务模型及其与镜像的关联信息。
package models

import "gorm.io/gorm"

// 服务表
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
