package grpc_service

// File: service/grpc_service/register.go
// Description: 节点注册

import (
	"context"
	"errors"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
)

// 处理节点注册请求
func (NodeService) Register(ctx context.Context, request *node_rpc.RegisterRequest) (pd *node_rpc.BaseResponse, err error) {
	// 初始化响应对象
	pd = new(node_rpc.BaseResponse)

	// 获取节点唯一标识（UID）
	uid := request.NodeUid

	// 检查节点是否已存在（通过UID查询）
	var model models.NodeModel
	dbErr := global.DB.Take(&model, "uid = ?", uid).Error
	if dbErr != nil {
		// 节点不存在，创建新节点记录
		model = models.NodeModel{
			Title:  request.SystemInfo.HostName, // 节点名称
			Uid:    uid,                         // 节点唯一标识
			IP:     request.Ip,                  // 节点IP地址
			Mac:    request.Mac,                 // 节点MAC地址
			Status: 1,                           // 节点状态：正常（1）
			SystemInfo: models.NodeSystemInfo{
				NodeVersion:         request.Version,                        // 节点版本号
				NodeCommit:          request.Commit,                         // 节点提交哈希
				HostName:            request.SystemInfo.HostName,            // 节点主机名
				DistributionVersion: request.SystemInfo.DistributionVersion, // 节点发行版版本号
				CoreVersion:         request.SystemInfo.CoreVersion,         // 节点内核版本号
				SystemType:          request.SystemInfo.SystemType,          // 节点系统类型
				StartTime:           request.SystemInfo.StartTime,           // 节点启动时间
			},
		}

		// 保存新节点到数据库
		dbErr = global.DB.Create(&model).Error
		if dbErr != nil {
			logrus.Errorf("节点创建失败: %v", dbErr)
			return nil, errors.New("节点创建失败")
		}
		logrus.Infof("节点注册成功（新节点），UID: %s", uid)
	} else {
		// 节点已存在，检查状态是否为正常（1）
		if model.Status != 1 {
			// 将状态更新为正常（1）（可能节点之前离线，重新注册时恢复在线状态）
			global.DB.Model(&model).Update("status", 1)
			logrus.Infof("节点注册成功（状态更新），UID: %s，旧状态: %d", uid, model.Status)
		} else {
			logrus.Infof("节点已在线，UID: %s", uid)
		}
	}

	return pd, nil
}
