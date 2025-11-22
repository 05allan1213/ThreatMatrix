package port_service

// File: service/port_service/load_tunnel.go
// Description: 提供端口转发配置的持久化加载功能，从数据库读取历史端口转发记录并自动启动对应的隧道服务

import (
	"honey_node/internal/global"
	"honey_node/internal/models"

	"github.com/sirupsen/logrus"
)

// LoadTunnel 从数据库加载历史端口转发配置并自动启动隧道服务
// 程序启动时调用，确保已配置的端口转发规则自动生效，实现配置持久化
func LoadTunnel() {
	var portList []models.PortModel
	// 从数据库查询所有已保存的端口转发记录
	global.DB.Find(&portList)
	logrus.Infof("加载端口转发记录 %d", len(portList))

	// 遍历端口转发记录，为每个配置启动独立的隧道服务
	for _, model := range portList {
		// 异步启动隧道，避免阻塞加载流程，支持多端口并发转发
		go Tunnel(model.LocalAddr, model.TargetAddr)
	}
}
