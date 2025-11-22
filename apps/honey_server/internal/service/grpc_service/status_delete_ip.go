package grpc_service

// File: service/grpc_service/status_delete_ip.go
// Description: 实现删除诱捕IP状态回调的gRPC接口处理逻辑，接收节点上报的删除结果并执行数据库批量删除操作

import (
	"context"
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
)

// StatusDeleteIP 删除诱捕IP状态回调的gRPC接口实现
func (NodeService) StatusDeleteIP(ctx context.Context, request *node_rpc.StatusDeleteIPRequest) (pd *node_rpc.BaseResponse, err error) {
	pd = new(node_rpc.BaseResponse) // 初始化gRPC响应结构体

	// 根据ID列表查询对应的诱捕IP记录
	var honeyIPList []models.HoneyIpModel
	global.DB.Find(&honeyIPList, "id in ?", request.HoneyIPIDList)

	// 记录删除回调的ID列表日志，便于追踪批量操作结果
	logrus.Infof("删除诱捕ip回调 %v", request.HoneyIPIDList)

	// 校验查询结果：若未找到任何记录则返回错误
	if len(honeyIPList) == 0 {
		return nil, fmt.Errorf("诱捕ip不存在 ")
	}

	// 执行批量删除操作（软删除/硬删除取决于模型配置的gorm标签）
	global.DB.Delete(&honeyIPList)

	return // 返回gRPC响应
}
