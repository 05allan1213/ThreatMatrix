package vs_net_api

// File: api/vs_net_api/enter.go
// Description: 虚拟网络配置 API，负责查询当前虚拟网络信息及更新虚拟网络配置

import (
	"fmt"
	"image_server/internal/config"
	"image_server/internal/core"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/cmd"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// VsNetApi 虚拟网络相关API的结构体
type VsNetApi struct {
}

// VsNetInfoView 获取虚拟网络配置信息
func (VsNetApi) VsNetInfoView(c *gin.Context) {
	res.OkWithData(global.Config.VsNet, c)
}

// VsNetRequest 更新虚拟网络配置的请求参数结构
type VsNetRequest struct {
	Name   string `json:"name" binding:"required"`   // 虚拟网络名称
	Prefix string `json:"prefix" binding:"required"` // 容器名称前缀
	Net    string `json:"net" binding:"required"`    // 子网配置（如10.2.0.0/24）
}

// VsNetUpdateView 更新虚拟网络配置
// 流程：1. 验证请求参数 2. 检查是否存在虚拟服务（存在则禁止修改）3. 删除旧网络 4. 创建新网络 5. 更新配置并保存
func (VsNetApi) VsNetUpdateView(c *gin.Context) {
	// 绑定并验证请求参数
	cr := middleware.GetBind[VsNetRequest](c)

	// 检查是否存在虚拟服务（有虚拟服务时禁止修改网络配置）
	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList)
	if len(serviceList) != 0 {
		res.FailWithMsg("存在虚拟服务，不可修改虚拟子网", c)
		return
	}

	// 删除原有虚拟网络（通过docker命令）
	command := fmt.Sprintf("docker network rm %s", global.Config.VsNet.Name)
	err := cmd.Cmd(command)
	if err != nil {
		logrus.Errorf("删除之前的虚拟网络失败 %s", err)
		res.FailWithMsg("删除之前的虚拟网络失败", c)
		return
	}

	// 创建新的虚拟网络（使用新配置的名称和子网）
	command = fmt.Sprintf("docker network create --driver bridge --subnet %s %s",
		cr.Net, cr.Name)
	err = cmd.Cmd(command)
	if err != nil {
		logrus.Errorf("创建虚拟网络失败 %s", err)
		res.FailWithMsg("创建虚拟网络失败", c)
		return
	}

	// 更新全局配置并保存到配置文件
	global.Config.VsNet = config.VsNet{
		Name:   cr.Name,
		Prefix: cr.Prefix,
		Net:    cr.Net,
	}
	core.SetConfig()

	res.OkWithMsg("修改虚拟网络成功", c)
}
