package command

// File: service/command/command_network_flush.go
// Description: 节点客户端中处理网卡刷新命令的逻辑实现，负责接收刷新请求、获取网卡列表并返回响应

import (
	"honey_node/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
)

// CmdNetworkFlush 处理网卡刷新命令
// 接收服务器的网卡刷新请求，根据过滤条件获取网卡列表，并将结果返回给服务器
func (nc *NodeClient) CmdNetworkFlush(request *node_rpc.CmdRequest) {
	logrus.Info("处理网卡刷新命令")

	// 提取请求中的网卡过滤条件（若存在）
	var filters []string
	if request.NetworkFlushInMessage != nil && len(request.NetworkFlushInMessage.FilterNetworkName) > 0 {
		filters = request.NetworkFlushInMessage.FilterNetworkName
	}

	// 根据过滤条件获取节点的网卡列表信息
	networkList, err := nc.getNetworkList(filters)
	if err != nil {
		logrus.Errorf("获取网卡列表失败: %v", err)
		return
	}

	// 构建网卡刷新响应，包含获取到的网卡列表
	response := &node_rpc.CmdResponse{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType, // 命令类型：网卡刷新
		TaskID:  request.TaskID,                       // 关联请求的任务ID，确保响应匹配
		NodeID:  nc.config.System.Uid,                 // 节点唯一标识
		NetworkFlushOutMessage: &node_rpc.NetworkFlushOutMessage{
			NetworkList: networkList, // 刷新后的网卡列表
		},
	}

	// 将响应发送到命令响应通道（非阻塞发送，避免通道满时阻塞）
	select {
	case nc.cmdResponseChan <- response:
		logrus.Debugf("已将响应加入发送队列: %+v", response)
	case <-nc.ctx.Done():
		// 上下文已取消（如节点退出），丢弃响应
		logrus.Warn("上下文已取消，丢弃响应")
	}
}
