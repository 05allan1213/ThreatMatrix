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
	"net"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VsCreateRequest 虚拟服务创建请求结构体
type VsCreateRequest struct {
	ImageID uint `json:"imageID" binding:"required"` // 指定的镜像ID
}

// 可用 IP 最大范围（支持 10.2.0.2 ~ 10.2.0.254）
const (
	maxIP = 254 // 最大可用IP 10.2.0.254
)

// getNextAvailableIP 获取下一个可用IP地址
// 逻辑说明：
// 1. 从配置中解析网段（如 10.2.0.0/24）
// 2. 查询数据库中目前已分配的最大 IP
// 3. 若无记录 → 返回网段基础地址 + 2（10.2.0.2）
// 4. 若有记录 → 最后一段 +1
// 5. 若超过 254 → 返回池已满错误
func getNextAvailableIP() (string, error) {
	// 从配置解析网段，例如 global.Config.VsNet.Net = "10.2.0.0/24"
	ip, _, err := net.ParseCIDR(global.Config.VsNet.Net)
	if err != nil {
		return "", err
	}
	ip4 := ip.To4() // 转换为 IPv4 格式

	var service models.ServiceModel

	// 查询服务表中 IP 最大的一条记录，用于推算下一个 IP
	err = global.DB.Order("ip DESC").First(&service).Error
	if err != nil {
		// 若数据库没有任何已使用 IP，则从网段 +2 开始分配
		if err.Error() == "record not found" {
			ip4[3] += 2 // 10.2.0.2
			return ip4.String(), nil
		}
		return "", fmt.Errorf("查询最大IP失败: %w", err)
	}

	// 将数据库中存储的 IP 解析为 net.IP
	serviceIP := net.ParseIP(service.IP)
	if serviceIP == nil {
		return "", fmt.Errorf("服务ip解析错误")
	}
	serviceIP4 := serviceIP.To4()

	// 检查是否超过最大分配范围
	if serviceIP4[3] >= maxIP {
		return "", fmt.Errorf("IP地址池已满")
	}

	// 正常分配：在当前最大 IP 的最后一段 +1
	newLastOctet := serviceIP4[3] + 1
	ip4[3] = newLastOctet
	return ip4.String(), nil
}

// VsCreateView 创建虚拟服务接口
// 详细流程：
// 1. 校验镜像 ID 是否存在、镜像状态是否可用
// 2. 确保 docker network 存在（没有则创建）
// 3. 检查该镜像是否已经运行过一个虚拟服务（防止重复创建）
// 4. 分配一个可用的 IP 地址
// 5. 通过 docker_service.RunContainer 运行容器
// 6. 将容器信息写入数据库 ServiceModel
func (VsApi) VsCreateView(c *gin.Context) {
	// 参数绑定与校验
	cr := middleware.GetBind[VsCreateRequest](c)

	// 查询镜像信息
	var image models.ImageModel
	err := global.DB.Take(&image, cr.ImageID).Error
	if err != nil {
		res.FailWithMsg("镜像不存在", c)
		return
	}

	// 判断镜像状态是否不可用
	if image.Status == 2 {
		res.FailWithMsg("镜像不可用", c)
		return
	}

	// 创建 docker 网络（若已存在则忽略）
	networkCommand := "docker network create --driver bridge --subnet 10.2.0.0/24 honey-hy >/dev/null 2>&1 || true"
	err = cmd.Cmd(networkCommand)
	if err != nil {
		logrus.Errorf("检查或创建Docker网络失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}

	// 判断该镜像是否已经运行容器（一个镜像只能创建一个虚拟服务）
	var service models.ServiceModel
	err = global.DB.Take(&service, "image_id = ?", cr.ImageID).Error
	if err == nil {
		res.FailWithMsg("此镜像已运行虚拟服务", c)
		return
	}

	// 分配下一可用 IP
	ip, err := getNextAvailableIP()
	if err != nil {
		logrus.Errorf("获取可用IP失败: %s", err)
		res.FailWithMsg("IP地址池已满，无法创建新服务", c)
		return
	}

	fmt.Println(ip) // 打印调试信息

	// 读取 docker 网络名称、容器前缀（来自配置）
	networkName := global.Config.VsNet.Name
	containerName := global.Config.VsNet.Prefix + image.ImageName

	// 实际启动容器
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

	// 输出用于调试的命令（实际执行 run 的是 RunContainer）
	command := fmt.Sprintf(
		"docker run -d --network %s --ip %s --name %s %s:%s",
		networkName, ip, containerName, image.ImageName, image.Tag,
	)
	fmt.Println(command)

	// 组装数据库记录
	var model = models.ServiceModel{
		Title:         image.Title,     // 服务标题
		ContainerName: containerName,   // 容器名称
		Agreement:     image.Agreement, // 协议类型（HTTP/HTTPS 等）
		ImageID:       image.ID,        // 镜像 ID
		IP:            ip,              // 分配 IP
		Port:          image.Port,      // 端口
		Status:        1,               // 状态 1=运行中
		ContainerID:   containerID,     // Docker 返回的容器 ID
	}

	// 写入数据库
	err = global.DB.Create(&model).Error
	if err != nil {
		logrus.Errorf("创建虚拟服务失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}

	res.OkWithMsg("创建虚拟服务成功", c)
}
