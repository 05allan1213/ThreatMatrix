package grpc_service

// File: service/grpc_service/node_service.go
// Description: 节点资源检测

import (
	"context"
	"errors"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
)

// 处理节点资源信息上报请求
func (NodeService) NodeResource(ctx context.Context, request *node_rpc.NodeResourceRequest) (pd *node_rpc.BaseResponse, err error) {
	// 初始化响应对象
	pd = new(node_rpc.BaseResponse)

	// 获取节点唯一标识（UID）
	uid := request.NodeUid

	// 检查节点是否存在（通过UID查询）
	var model models.NodeModel
	dbErr := global.DB.Take(&model, "uid = ?", uid).Error
	if dbErr != nil {
		return nil, errors.New("节点不存在")
	}

	// 构建需要更新的节点资源信息（仅包含资源相关字段）
	newModel := models.NodeModel{
		Resource: models.NodeResource{
			CpuCount:           int(request.ResourceInfo.CpuCount),         // CPU核心数
			CpuUseRate:         float64(request.ResourceInfo.CpuUseRate),   // CPU使用率
			MemTotal:           request.ResourceInfo.MemTotal,              // 总内存
			MemUseRate:         float64(request.ResourceInfo.MemUseRate),   // 内存使用率
			DiskTotal:          request.ResourceInfo.DiskTotal,             // 总磁盘空间
			DiskUseRate:        float64(request.ResourceInfo.DiskUseRate),  // 磁盘使用率
			NodePath:           request.ResourceInfo.NodePath,              // 节点路径
			NodeResourceOccupy: request.ResourceInfo.NodeResourceOccupancy, // 节点资源占用率
		},
	}

	// 更新数据库中节点的资源信息
	dbErr = global.DB.Model(&model).Updates(newModel).Error
	if dbErr != nil {
		logrus.Errorf("节点资源状态更新失败: %v", dbErr)
		return nil, errors.New("节点资源状态更新失败")
	}

	logrus.Infof("节点资源信息更新成功，UID: %s", uid)
	return pd, nil
}
