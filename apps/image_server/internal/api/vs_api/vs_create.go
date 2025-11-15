package vs_api

// File: api/vs_api/vs_create.go
// Description: 虚拟服务 API，负责虚拟服务创建、容器运行、IP 分配以及状态检测。

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/docker_service"
	"image_server/internal/utils/res"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VsCreateRequest 创建虚拟服务的请求参数结构
type VsCreateRequest struct {
	ImageID uint `json:"imageID" binding:"required"` // 需创建虚拟服务的镜像ID，必填项
}

// 基础IP地址与子网配置常量
const (
	maxIP = 254 // 可用IP地址范围的最大值（对应最后一段IP：10.2.0.254）
)

// getNextAvailableIP 获取下一个可用的IP地址
// 逻辑：基于配置的基础网段，从数据库查询已分配的最大IP，递增生成下一个IP；若未分配过则从10.2.0.2开始
func getNextAvailableIP() (string, error) {
	// 解析配置中的基础网段（如10.2.0.0/24）
	ip, _, err := net.ParseCIDR(global.Config.VsNet.Net)
	if err != nil {
		return "", err
	}
	ip4 := ip.To4() // 转换为IPv4格式

	// 查询数据库中当前已分配的最大IP记录
	var service models.ServiceModel
	err = global.DB.Order("ip DESC").First(&service).Error
	if err != nil {
		// 若没有任何分配记录，从.2开始分配（10.2.0.2）
		if err.Error() == "record not found" {
			ip4[3] = 2
			return ip4.String(), nil
		}
		return "", fmt.Errorf("查询最大IP失败: %w", err)
	}

	// 解析数据库中已存在的IP地址
	serviceIP := net.ParseIP(service.IP)
	if serviceIP == nil {
		return "", fmt.Errorf("服务ip解析错误")
	}
	serviceIP4 := serviceIP.To4()

	// 检查IP是否已达最大可用范围（10.2.0.254）
	if serviceIP4[3] >= maxIP {
		return "", fmt.Errorf("IP地址池已满")
	}

	// 生成下一个可用IP（最后一段+1）
	newLastOctet := serviceIP4[3] + 1
	ip4[3] = newLastOctet
	return ip4.String(), nil
}

// VsCreateView 创建虚拟服务的API入口函数
func (VsApi) VsCreateView(c *gin.Context) {
	// 绑定并验证请求参数（从上下文获取VsCreateRequest结构体）
	cr := middleware.GetBind[VsCreateRequest](c)

	// 查询镜像信息（根据请求的ImageID）
	var image models.ImageModel
	err := global.DB.Take(&image, cr.ImageID).Error
	if err != nil {
		res.FailWithMsg("镜像不存在", c)
		return
	}
	// 检查镜像状态是否可用（状态2为不可用）
	if image.Status == 2 {
		res.FailWithMsg("镜像不可用", c)
		return
	}

	// 检查该镜像是否已创建过虚拟服务（避免重复创建）
	var service models.ServiceModel
	err = global.DB.Take(&service, "image_id = ?", cr.ImageID).Error
	if err == nil {
		res.FailWithMsg("此镜像已运行虚拟服务", c)
		return
	}

	// 分配可用IP地址
	ip, err := getNextAvailableIP()
	if err != nil {
		logrus.Errorf("获取可用IP失败: %s", err)
		res.FailWithMsg("IP地址池已满，无法创建新服务", c)
		return
	}

	fmt.Println(ip) // 打印分配的IP（调试用）

	// 运行容器（基于镜像信息、网络配置、分配的IP）
	networkName := global.Config.VsNet.Name                       // 网络名称（从配置获取）
	containerName := global.Config.VsNet.Prefix + image.ImageName // 容器名称（前缀+镜像名）
	// 调用docker服务创建容器，返回容器ID
	containerID, err := docker_service.RunContainer(containerName, networkName, ip, fmt.Sprintf("%s:%s", image.ImageName, image.Tag))
	if err != nil {
		logrus.Errorf("创建虚拟服务失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}

	// 打印docker命令（调试用，便于手动验证）
	command := fmt.Sprintf("docker run -d --network %s --ip %s --name %s %s:%s",
		networkName, ip, containerName, image.ImageName, image.Tag)
	fmt.Println(command)

	// 在数据库中创建虚拟服务记录
	var model = models.ServiceModel{
		Title:         image.Title,     // 服务标题（复用镜像标题）
		ContainerName: containerName,   // 容器名称
		Agreement:     image.Agreement, // 协议类型（复用镜像配置）
		ImageID:       image.ID,        // 关联的镜像ID
		IP:            ip,              // 分配的IP地址
		Port:          image.Port,      // 端口（复用镜像配置）
		Status:        1,               // 初始状态：正常（1）
		ContainerID:   containerID,     // 容器ID
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		logrus.Errorf("创建虚拟服务失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}

	// 启动后台协程，定时检测容器状态（分多个时间点检测）
	go func(model *models.ServiceModel) {
		// 检测时间点：5秒、20秒、1分钟、5分钟、1小时后
		var delayList = []<-chan time.Time{
			time.After(5 * time.Second),
			time.After(20 * time.Second),
			time.After(1 * time.Minute),
			time.After(5 * time.Minute),
			time.After(1 * time.Hour),
		}
		// 依次等待每个时间点，执行状态检测
		for _, times := range delayList {
			<-times
			ContainerStatus(model)
		}
	}(&model)

	// 返回创建成功响应
	res.OkWithMsg("创建虚拟服务成功", c)
}

// ContainerStatus 检测容器运行状态并更新数据库记录
// 逻辑：对比容器实际状态与数据库记录状态，若不一致则更新数据库
func ContainerStatus(model *models.ServiceModel) {
	logrus.Infof("检测容器状态 %s", model.ContainerName) // 记录检测日志
	var newModel models.ServiceModel               // 用于存储待更新的状态信息

	// 获取容器状态（通过docker服务查询）
	containers, err := docker_service.PrefixContainerStatus(model.ContainerName)
	var isUpdate bool // 是否需要更新数据库
	var state string  // 容器当前状态描述

	// 处理查询错误（如docker服务异常）
	if err != nil {
		newModel.Status = 2             // 状态：异常（2）
		newModel.ErrorMsg = err.Error() // 记录错误信息
		isUpdate = true
		state = err.Error()
	}

	// 处理容器不存在的情况（查询结果数量不等于1）
	if len(containers) != 1 {
		newModel.Status = 2
		newModel.ErrorMsg = "容器不存在"
		isUpdate = true
		state = newModel.ErrorMsg
	} else {
		container := containers[0] // 获取唯一容器信息

		// 容器运行中，但数据库记录为非正常状态 → 更新为正常
		if container.State == "running" && model.Status != 1 {
			newModel.Status = 1
			newModel.ErrorMsg = ""
			isUpdate = true
			state = container.State
		}

		// 容器非运行中，但数据库记录为正常状态 → 更新为异常
		if container.State != "running" && model.Status == 1 {
			newModel.Status = 2
			newModel.ErrorMsg = fmt.Sprintf("%s(%s)", container.State, container.Status) // 记录状态详情
			isUpdate = true
			state = container.State
		}
	}

	// 若状态有变化，更新数据库并记录日志
	if isUpdate {
		logrus.Infof("%s 容器存在状态修改 %s => %s", model.ContainerName, model.State(), state)
		global.DB.Model(model).Updates(map[string]any{
			"status":    newModel.Status,
			"error_msg": newModel.ErrorMsg,
		})
	}
}
