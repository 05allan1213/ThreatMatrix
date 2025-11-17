package net_api

// File: api/net_api/scan.go
// Description: 网络扫描API

import (
	"context"
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ScanView 处理网络扫描的请求
func (NetApi) ScanView(c *gin.Context) {
	// 获取并绑定请求参数（包含要扫描的网络ID）
	cr := middleware.GetBind[models.IDRequest](c)

	// 查询指定ID的网络信息，并预加载关联的节点信息
	var model models.NetModel
	if err := global.DB.Preload("NodeModel").Take(&model, cr.Id).Error; err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}

	// 检查网络所属节点是否处于运行状态（状态1为运行）
	if model.NodeModel.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 通过节点唯一标识（Uid）获取节点命令通道（用于发送和接收grpc命令）
	cmd, ok := grpc_service.GetNodeCommand(model.NodeModel.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 构建网络扫描请求参数
	req := &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetScanType,                  // 命令类型：网络扫描
		TaskID:  fmt.Sprintf("netScan-%d", time.Now().UnixNano()), // 生成唯一任务ID（基于时间戳）
		NetScanInMessage: &node_rpc.NetScanInMessage{
			Network:      model.Network,            // 目标网络地址
			IpRange:      model.CanUseHoneyIPRange, // 可用蜜罐IP范围
			FilterIPList: []string{},               // 过滤IP列表（暂为空）
			NetID:        uint32(model.ID),         // 网络ID
		},
	}

	// 创建带30秒超时的上下文，防止请求无响应时阻塞
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// 发送扫描请求到节点的命令通道
	select {
	case cmd.ReqChan <- req:
		logrus.Debugf("已向节点 %s 发送扫描请求", model.NodeModel.Uid)
	case <-ctx.Done():
		// 发送请求超时
		res.FailWithMsg("发送命令超时", c)
		return
	}

	// 等待节点返回扫描结果（循环监听响应通道）
label:
	for {
		select {
		case response := <-cmd.ResChan:
			// 接收到节点的响应
			logrus.Debugf("已接收节点 %s 的扫描响应", model.NodeModel.Uid)
			message := response.NetScanOutMessage
			fmt.Printf("网络扫描 数据 %v\n", message) // 打印扫描数据
			if message.ErrMsg != "" {
				res.FailWithMsg("扫描错误"+message.ErrMsg, c)
				return
			}
			if message.End {
				// 扫描结束，跳出循环
				break label
			}
		case <-ctx.Done():
			// 等待响应超时
			res.FailWithMsg("获取响应超时", c)
			return
		}
	}

	// 扫描成功，返回提示信息
	res.OkWithMsg("扫描成功", c)
}
