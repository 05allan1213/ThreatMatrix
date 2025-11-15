package cron_service

// File: service/cron_service/vs_health.go
// Description: 定时检查 Docker 容器健康状态，同步服务状态至数据库

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/models"
	"image_server/internal/service/docker_service"

	"github.com/sirupsen/logrus"
)

// VsHealth 定时任务：检测所有容器运行状态并同步至数据库
func VsHealth() {
	logrus.Infof("获取前缀 %s 的容器状态", global.Config.VsNet.Prefix)
	allContainers, err := docker_service.PrefixContainerStatus(global.Config.VsNet.Prefix)
	if err != nil {
		logrus.Errorf("容器状态检测失败 %s", err)
		return
	}

	var list []models.ServiceModel
	global.DB.Find(&list)

	// 构建 containerID 与服务模型的映射
	var containerMap = map[string]*models.ServiceModel{}
	for _, model := range list {
		containerMap[model.ContainerID] = &model
	}

	// 遍历 Docker 容器列表
	for _, container := range allContainers {
		containerID := container.ID[:12]
		model, ok := containerMap[containerID]
		if !ok {
			continue
		}

		var newModel models.ServiceModel
		var isUpdate bool

		// 状态：数据库认为“不正常”，但 Docker 实际“正常”
		if container.State == "running" && model.Status != 1 {
			newModel.Status = 1
			newModel.ErrorMsg = ""
			isUpdate = true
		}

		// 状态：数据库认为“正常”，但 Docker 实际“不正常”
		if container.State != "running" && model.Status == 1 {
			newModel.Status = 2
			newModel.ErrorMsg = fmt.Sprintf("%s(%s)", container.State, container.Status)
			isUpdate = true
		}

		// 有状态变更才写库
		if isUpdate {
			logrus.Infof("%s 容器存在状态修改 %s => %s", model.ContainerName, model.State(), container.State)
			global.DB.Model(model).Updates(map[string]any{
				"status":    newModel.Status,
				"error_msg": newModel.ErrorMsg,
			})
		}
	}
}
