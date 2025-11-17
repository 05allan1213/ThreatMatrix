package node_network_api

// File: api/node_network_api/flush.go
// Description: 提供带超时控制的节点网卡信息刷新API

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// 处理指定节点的网卡信息刷新请求（带超时控制）
func (NodeNetworkApi) FlushView(c *gin.Context) {
	// 获取并绑定请求参数（节点ID），通过中间件处理参数校验
	cr := middleware.GetBind[models.IDRequest](c)

	// 根据ID查询节点信息，验证节点是否存在
	var model models.NodeModel
	if err := global.DB.Take(&model, cr.Id).Error; err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}

	// 验证节点是否处于运行状态（状态1表示运行中）
	if model.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 通过封装的方法获取节点对应的命令交互实例，验证节点是否在线
	cmd, ok := grpc_service.GetNodeCommand(model.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 构建网卡刷新命令请求，使用当前时间戳的纳秒数作为唯一任务ID
	req := &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,           // 命令类型：网卡刷新
		TaskID:  fmt.Sprintf("flush-%d", time.Now().UnixNano()), // 任务ID，确保唯一性
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"}, // 过滤"hy-"前缀的网卡
		},
	}

	// 创建带30秒超时的上下文，控制命令发送和响应接收的超时时间
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel() // 确保上下文最终被取消，释放资源

	// 向节点的命令请求通道发送刷新命令，带超时控制
	select {
	case cmd.ReqChan <- req:
		// 命令发送成功，记录调试日志
		logrus.Debugf("已向节点[%s]发送网卡刷新请求", model.Uid)
	case <-ctx.Done():
		// 命令发送超时，返回错误提示
		res.FailWithMsg("发送命令超时", c)
		return
	}

	// 等待节点返回的刷新结果，带超时控制
	select {
	case response := <-cmd.ResChan:
		// 成功接收响应，记录调试日志并返回结果
		logrus.Debugf("已接收节点[%s]的网卡刷新响应", model.Uid)
		res.OkWithData(response.NetworkFlushOutMessage, c)
	case <-ctx.Done():
		// 接收响应超时，返回错误提示
		res.FailWithMsg("获取响应超时", c)
		return
	}
}
