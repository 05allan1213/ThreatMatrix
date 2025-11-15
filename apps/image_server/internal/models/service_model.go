package models

import (
	"errors"
	"fmt"
	"image_server/internal/utils/cmd"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// File: models/service_model.go
// Description: 定义诱捕服务的数据模型及其与镜像、节点的关联关系。

// 服务模型
type ServiceModel struct {
	Model
	Title         string     `json:"title"`                       // 服务名称 用镜像名称
	Agreement     int8       `json:"agreement"`                   // 协议
	ImageID       uint       `json:"imageID"`                     // 镜像ID
	ImageModel    ImageModel `gorm:"foreignKey:ImageID" json:"-"` // 关联镜像模型
	IP            string     `json:"ip"`                          // 容器IP
	Port          int        `json:"port"`                        // 容器端口
	Status        int8       `json:"status"`                      // 容器状态
	ErrorMsg      string     `json:"errorMsg"`                    // 错误信息
	HoneyIPCount  int        `json:"honeyIPCount"`                // 诱捕IP数量
	ContainerID   string     `json:"containerID"`                 // 容器ID
	ContainerName string     `json:"containerName"`               // 容器名称
}

// 获取服务状态
func (s *ServiceModel) State() string {
	switch s.Status {
	case 1:
		return "running"
	}
	return "error"
}

func (s *ServiceModel) BeforeDelete(tx *gorm.DB) error {
	// 判断是否有关联的端口转发
	var count int64
	tx.Model(HoneyPortModel{}).Where("service_id = ?", s.ID).Count(&count)
	if count > 0 {
		return errors.New("存在端口转发，不能删除虚拟服务")
	}

	command := fmt.Sprintf("docker rm -f %s", s.ContainerName)
	err := cmd.Cmd(command)
	if err != nil {
		logrus.Errorf("删除容器失败 %s", err)
		return err
	}
	return nil
}
