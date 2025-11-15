package models

import (
	"fmt"
	"image_server/internal/utils/cmd"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// File: models/image_model.go
// Description: 定义镜像资源的数据模型及其与服务的关联关系。

// 镜像模型
type ImageModel struct {
	Model
	ImageName     string         `gorm:"size:64" json:"imageName"`     // 镜像名称
	Title         string         `gorm:"size:64" json:"title"`         // 镜像别名
	Port          int            `json:"port"`                         // 镜像端口
	DockerImageID string         `gorm:"size:32" json:"dockerImageID"` // Docker镜像ID
	ServiceList   []ServiceModel `gorm:"foreignKey:ImageID" json:"-"`  // 关联的虚拟服务列表
	Tag           string         `gorm:"size:32" json:"tag"`           // 镜像标签
	Agreement     int8           `json:"agreement"`                    // 镜像协议
	ImagePath     string         `gorm:"size:256" json:"-"`            // 镜像文件
	Status        int8           `json:"status"`                       // 镜像状态 1 成功
	Logo          string         `gorm:"size:256" json:"logo"`         // 镜像logo
	Desc          string         `gorm:"size:512" json:"desc"`         // 镜像描述
}

// BeforeDelete 删除镜像之前执行, 清理docker镜像和镜像文件
func (i *ImageModel) BeforeDelete(tx *gorm.DB) error {
	// 删除docker镜像
	command := fmt.Sprintf("docker rmi %s", i.DockerImageID)
	err := cmd.Cmd(command)
	if err != nil {
		return err
	}
	// 删除镜像文件
	logrus.Infof("删除镜像文件 %s", i.ImagePath)
	err = os.Remove(i.ImagePath)
	if err != nil {
		return err
	}
	return nil
}
