package vs_api

// File: api/vs_api/vs_create.go
// Description: 虚拟服务创建接口，根据指定镜像运行容器并写入数据库记录

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/cmd"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VsCreateRequest 虚拟服务创建请求结构体
type VsCreateRequest struct {
	ImageID uint `json:"imageID" binding:"required"` // 指定的镜像ID
}

// VsCreateView 创建虚拟服务接口
// 1. 校验镜像是否存在且可用
// 2. 运行 Docker 命令创建容器
// 3. 写入 ServiceModel 数据库记录
func (VsApi) VsCreateView(c *gin.Context) {
	cr := middleware.GetBind[VsCreateRequest](c)

	// 查询镜像信息
	var image models.ImageModel
	err := global.DB.Take(&image, cr.ImageID).Error
	if err != nil {
		res.FailWithMsg("镜像不存在", c)
		return
	}
	if image.Status == 2 {
		res.FailWithMsg("镜像不可用", c)
		return
	}

	// 确保Docker网络存在
	networkCommand := "docker network create --driver bridge --subnet 10.2.0.0/24 honey-hy >/dev/null 2>&1 || true"
	err = cmd.Cmd(networkCommand)
	if err != nil {
		logrus.Errorf("检查或创建Docker网络失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
	}

	// 使用 docker 命令运行容器
	// 示例命令:
	// docker network create --driver bridge --subnet 10.2.0.0/24 honey-hy
	// docker run -d --network honey-hy --ip 10.2.0.10 --name my_container image_name:tag
	ip := "10.2.0.10"
	command := fmt.Sprintf(
		"docker run -d --network honey-hy --ip %s --name %s_%d %s:%s",
		ip, image.ImageName, cr.ImageID, image.ImageName, image.Tag,
	)

	// 构建 ServiceModel 数据结构
	var model = models.ServiceModel{
		Title:       image.ImageName,
		Agreement:   image.Agreement,
		ImageID:     image.ID,
		IP:          ip,
		Port:        image.Port,
		Status:      1,
		ContainerID: "",
	}

	// 执行 Docker 命令
	err = cmd.Cmd(command)
	if err != nil {
		logrus.Errorf("创建虚拟服务失败 %s", err)
		res.FailWithMsg("创建虚拟服务失败", c)
		return
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
