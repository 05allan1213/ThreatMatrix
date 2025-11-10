// Package models 定义诱捕服务所使用的数据实体。
//
// 本文件描述镜像资源的元数据结构，用于支撑服务部署。
package models

import "gorm.io/gorm"

// 镜像模型
type ImageModel struct {
	gorm.Model
	ImageName string `json:"imageName"` // 镜像名称
	Title     string `json:"title"`     // 对外展示名称
	Port      int    `json:"port"`      // 镜像端口
	ImageID   string `json:"imageID"`   // 镜像ID
	Tag       string `json:"tag"`       // 镜像标签
	Agreement int8   `json:"agreement"` // 镜像协议
	ImagePath string `json:"imagePath"` // 镜像文件
	Status    int8   `json:"status"`    // 镜像状态
	Logo      string `json:"logo"`      // 镜像的logo
	Desc      string `json:"desc"`      // 镜像描述
}
