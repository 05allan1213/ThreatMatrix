package grpc_service

// File: service/grpc_service/register.go
// Description: 实现grpc服务的节点注册接口

import (
	"context"
	"errors"
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
)

// 处理节点注册的grpc接口实现
func (NodeService) Register(ctx context.Context, request *node_rpc.RegisterRequest) (pd *node_rpc.BaseResponse, err error) {
	// 初始化响应对象
	pd = new(node_rpc.BaseResponse)

	// 从请求中获取节点唯一标识（UID）
	uid := request.NodeUid

	// 声明节点模型变量，用于数据库操作
	var model models.NodeModel

	// 检查数据库中是否已存在该UID的节点
	err1 := global.DB.Take(&model, "uid = ?", uid).Error
	if err1 != nil {
		// 节点不存在，创建新节点记录
		model = models.NodeModel{
			Title:  request.SystemInfo.HostName, // 节点标题默认使用主机名
			Uid:    uid,                         // 节点唯一标识
			IP:     request.Ip,                  // 节点IP地址
			Mac:    request.Mac,                 // 节点MAC地址
			Status: 1,                           // 节点状态（1表示正常）
			SystemInfo: models.NodeSystemInfo{ // 系统信息字段赋值
				NodeVersion:         request.Version,
				NodeCommit:          request.Commit,
				HostName:            request.SystemInfo.HostName,
				DistributionVersion: request.SystemInfo.DistributionVersion,
				CoreVersion:         request.SystemInfo.CoreVersion,
				SystemType:          request.SystemInfo.SystemType,
				StartTime:           request.SystemInfo.StartTime,
			},
		}

		// 保存新节点到数据库
		err1 = global.DB.Create(&model).Error
		if err1 != nil {
			logrus.Errorf("节点创建失败 %s", err)
			return nil, errors.New("节点创建失败")
		}

		// 处理节点网卡信息，构建网卡记录列表
		var networkList []models.NodeNetworkModel
		for _, message := range request.NetworkList {
			networkList = append(networkList, models.NodeNetworkModel{
				NodeID:  model.ID,           // 关联的节点ID
				Network: message.Network,    // 网卡名称
				IP:      message.Ip,         // 网卡IP地址
				Mask:    fmt.Sprintf("%d", message.Mask), // 子网掩码长度
				Status:  2,                  // 网卡状态（2表示未启用）
			})
		}

		// 若存在网卡信息，批量保存到数据库
		if len(networkList) > 0 {
			err = global.DB.Create(&networkList).Error
			if err != nil {
				logrus.Errorf("节点网卡保存失败 %s", err)
				return nil, errors.New("节点网卡保存失败")
			}
		}
	}

	// 若节点已存在且状态不为在线（1），更新状态为在线
	if model.Status != 1 {
		global.DB.Model(&model).Update("status", 1)
	}

	return
}
