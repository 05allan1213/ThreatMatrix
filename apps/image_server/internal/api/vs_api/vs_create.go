package vs_api

// File: api/vs_api/vs_create.go
// Description: 虚拟服务创建接口，根据指定镜像运行容器并写入数据库记录

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/docker_service"
	"image_server/internal/utils/cmd"
	"image_server/internal/utils/res"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VsCreateRequest 虚拟服务创建请求结构体
type VsCreateRequest struct {
	ImageID uint `json:"imageID" binding:"required"` // 指定的镜像ID
}

// 基础IP地址和子网掩码
const (
	baseIP      = "10.2.0.0" // 网络段基础地址
	netmask     = 24         // 子网掩码
	startIP     = 2          // 从10.2.0.2开始分配
	maxIP       = 254        // 最大可用IP 10.2.0.254
	reservedIPs = 1          // 预留IP数量
)

// getNextAvailableIP 获取下一个可用IP地址
// 逻辑说明：
// 1. 查询 service 表中最大的 IP
// 2. 如果没有记录 → 返回起始IP 10.2.0.2
// 3. 如果有记录 → 将最后一段 +1 得到新IP
// 4. 若超过最大范围，则返回错误
func getNextAvailableIP() (string, error) {
	var service models.ServiceModel

	// 从数据库中查询 IP 最大的记录，用于计算下一个IP
	err := global.DB.Order("ip DESC").First(&service).Error
	if err != nil {
		// 数据表中还没有任何服务记录时，从起始IP分配
		if err.Error() == "record not found" {
			return "10.2.0.2", nil
		}
		return "", fmt.Errorf("查询最大IP失败: %w", err)
	}

	// 将 IP 按 "." 分割为 4 段
	ipParts := strings.Split(service.IP, ".")
	if len(ipParts) != 4 {
		return "", fmt.Errorf("无效的IP格式: %s", service.IP)
	}

	// 获取最后一段并递增
	lastOctet, err := strconv.Atoi(ipParts[3])
	if err != nil {
		return "", fmt.Errorf("解析IP最后一段失败: %s", service.IP)
	}

	// 检查可用范围
	if lastOctet >= maxIP {
		return "", fmt.Errorf("IP地址池已满")
	}

	// 构造新的IP地址
	newLastOctet := lastOctet + 1
	newIP := fmt.Sprintf("10.2.0.%d", newLastOctet)
	return newIP, nil
}

// VsCreateView 虚拟服务创建接口
// 实现步骤：
// 1. 校验镜像是否存在且可用
// 2. 确保 docker 网络已创建
// 3. 校验该镜像是否已运行服务（防重复）
// 4. 分配新的IP地址
// 5. 调用 docker_service 启动容器
// 6. 将服务记录写入数据库
func (VsApi) VsCreateView(c *gin.Context) {
	// 解析并绑定请求体
	cr := middleware.GetBind[VsCreateRequest](c)

	// 根据 ImageID 查询镜像记录
	var image models.ImageModel
	err := global.DB.Take(&image, cr.ImageID).Error
	if err != nil {
		res.FailWithMsg("镜像不存在", c)
		return
	}

	// 镜像状态为 2 表示不可用
	if image.Status == 2 {
		res.FailWithMsg("镜像不可用", c)
		return
	}

	// 确保 docker 网络存在
	// 如果网络 honey-hy 不存在，则创建；存在则忽略错误
	networkCommand := "docker network create --driver bridge --subnet 10.2.0.0/24 honey-hy >/dev/null 2>&1 || true"
	err = cmd.Cmd(networkCommand)
	if err != nil {
		logrus.Errorf("检查或创建Docker网络失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}

	// 判断此镜像是否已经创建过服务，避免重复创建
	var service models.ServiceModel
	err = global.DB.Take(&service, "image_id = ?", cr.ImageID).Error
	if err == nil {
		res.FailWithMsg("此镜像已运行虚拟服务", c)
		return
	}

	// 获取新的可分配IP
	ip, err := getNextAvailableIP()
	if err != nil {
		logrus.Errorf("获取可用IP失败: %s", err)
		res.FailWithMsg("IP地址池已满，无法创建新服务", c)
		return
	}

	fmt.Println(ip)

	// 构建容器名称（加上 hy_ 前缀）
	networkName := "honey-hy"
	containerName := "hy_" + image.ImageName

	// 使用 docker_service.RunContainer 封装方法运行容器
	containerID, err := docker_service.RunContainer(
		containerName,
		networkName,
		ip,
		fmt.Sprintf("%s:%s", image.ImageName, image.Tag),
	)
	if err != nil {
		logrus.Errorf("创建虚拟服务失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}

	// 打印实际执行的 Docker 命令（用于调试）
	command := fmt.Sprintf(
		"docker run -d --network honey-hy --ip %s --name %s %s:%s",
		ip, image.ImageName, image.ImageName, image.Tag,
	)
	fmt.Println(command)

	// 组装 ServiceModel 记录
	var model = models.ServiceModel{
		Title:         image.Title,     // 服务名称
		ContainerName: containerName,   // 容器名称
		Agreement:     image.Agreement, // 协议类型（例如 http/https）
		ImageID:       image.ID,        // 镜像ID
		IP:            ip,              // 分配的IP
		Port:          image.Port,      // 服务端口
		Status:        1,               // 状态：1-运行中
		ContainerID:   containerID,     // Docker容器ID
	}

	// 将新服务写入数据库
	err = global.DB.Create(&model).Error
	if err != nil {
		logrus.Errorf("创建虚拟服务失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}

	res.OkWithMsg("创建虚拟服务成功", c)
}
