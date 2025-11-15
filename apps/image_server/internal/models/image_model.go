package models

// File: models/image_model.go
// Description: 定义镜像资源的数据模型及其与服务的关联关系。

// 镜像模型
type ImageModel struct {
	Model
	ImageName     string `gorm:"size:64" json:"imageName"`     // 镜像名称
	Title         string `gorm:"size:64" json:"title"`         // 镜像别名
	Port          int    `json:"port"`                         // 镜像端口
	DockerImageID string `gorm:"size:32" json:"dockerImageID"` // 镜像ID
	Tag           string `gorm:"size:32" json:"tag"`           // 镜像标签
	Agreement     int8   `json:"agreement"`                    // 协议
	ImagePath     string `gorm:"size:256" json:"-"`            // 镜像文件
	Status        int8   `json:"status"`                       // 镜像状态 1 成功
	Logo          string `gorm:"size:256" json:"logo"`         // 镜像的logo
	Desc          string `gorm:"size:512" json:"desc"`         // 镜像描述
}
