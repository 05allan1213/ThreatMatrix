package cron_service

// File: service/cron_service/resource.go
// Description: 提供节点资源信息定时上报功能

import (
	"context"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/info"
	"os"

	"github.com/sirupsen/logrus"
)

// 定时上报节点资源信息到管理服务
func Resource() {
	// 检查grpc客户端是否已连接，未连接则记录错误并返回
	if global.GrpcClient == nil {
		logrus.Errorf("管理未连接，放弃上报")
		return
	}

	// 获取当前节点的工作目录路径
	nodePath, _ := os.Getwd()

	// 获取节点资源信息（CPU、内存、磁盘等使用情况）
	resourceInfo, err := info.GetResourceInfo(nodePath)
	if err != nil {
		logrus.Errorf("节点资源信息获取失败 %s", err)
		return
	}

	// 构建资源上报请求，填充节点唯一标识（UID）和资源信息
	req := node_rpc.NodeResourceRequest{
		NodeUid: global.Config.System.Uid,
		ResourceInfo: &node_rpc.ResourceMessage{
			CpuCount:              resourceInfo.CpuCount,
			CpuUseRate:            resourceInfo.CpuUseRate,
			MemTotal:              resourceInfo.MemTotal,
			MemUseRate:            resourceInfo.MemUseRate,
			DiskTotal:             resourceInfo.DiskTotal,
			DiskUseRate:           resourceInfo.DiskUseRate,
			NodePath:              resourceInfo.NodePath,
			NodeResourceOccupancy: resourceInfo.NodeResourceOccupancy,
		},
	}

	// 调用grpc客户端方法上报资源信息
	_, err = global.GrpcClient.NodeResource(context.Background(), &req)
	if err != nil {
		logrus.Errorf("节点资源信息上报失败 %s", err)
		return
	}

	// 上报成功，记录信息日志
	// logrus.Infof("节点资源信息上报成功 %v", req)
}
