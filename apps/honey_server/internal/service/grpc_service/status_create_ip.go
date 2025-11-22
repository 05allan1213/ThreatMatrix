package grpc_service

// File: service/grpc_service/status_create_ip.go
// Description: 实现创建诱捕IP状态上报的gRPC接口处理逻辑

import (
	"context"
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
)

// StatusCreateIP 诱捕IP创建状态上报的gRPC接口实现
func (NodeService) StatusCreateIP(ctx context.Context, request *node_rpc.StatusCreateIPRequest) (pd *node_rpc.BaseResponse, err error) {
	pd = new(node_rpc.BaseResponse) // 初始化gRPC响应结构体

	// 根据诱捕IP ID查询对应的记录
	var honeyIPModel models.HoneyIpModel
	err1 := global.DB.Take(&honeyIPModel, request.HoneyIPID).Error
	if err1 != nil {
		return nil, fmt.Errorf("诱捕ip不存在 %d", request.HoneyIPID)
	}

	// 设置状态：默认2表示创建成功，若存在错误信息则设为3（创建失败）
	var status int8 = 2
	if request.ErrMsg != "" {
		status = 3                                   // 状态3表示创建失败
		logrus.Errorf("创建诱捕ip失败 %s", request.ErrMsg) // 记录失败日志
	}

	// 更新诱捕IP记录的MAC地址、所属网络及状态
	global.DB.Model(&honeyIPModel).Updates(models.HoneyIpModel{
		Mac:     request.Mac,     // 节点上报的MAC地址
		Network: request.Network, // 节点上报的所属网络接口
		Status:  status,          // 最终创建状态
	})

	return // 返回gRPC响应
}
