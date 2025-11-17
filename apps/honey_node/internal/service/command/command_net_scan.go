package command

// File: service/command/net_scan.go
// Description: 节点客户端中处理网络扫描命令的逻辑实现，负责接收网络扫描请求并返回扫描进度及结果

import (
	"fmt"
	"honey_node/internal/rpc/node_rpc"
	"time"
)

// CmdNetScan 处理网络扫描命令
// 接收网络扫描请求，模拟扫描过程并返回阶段性进度及最终结果
func (nc *NodeClient) CmdNetScan(request *node_rpc.CmdRequest) {
	// 从请求中获取网络扫描的具体参数（如目标网络、IP范围等）
	req := request.GetNetScanInMessage()
	// 打印调试信息，输出网络扫描参数
	fmt.Printf("网络扫描 %v\n", req)

	// 发送第一阶段响应：扫描未结束，进度0%，并附带一个示例IP
	nc.cmdResponseChan <- &node_rpc.CmdResponse{
		CmdType: node_rpc.CmdType_cmdNetScanType, // 命令类型：网络扫描
		TaskID:  request.TaskID,                  // 关联的任务ID，与请求保持一致
		NodeID:  nc.config.System.Uid,            // 节点唯一标识
		NetScanOutMessage: &node_rpc.NetScanOutMessage{
			End:      false,           // 未结束
			Progress: 0,               // 进度0%
			Ip:       "192.168.100.1", // 示例扫描到的IP
		},
	}

	// 模拟扫描过程，等待2秒
	time.Sleep(2 * time.Second)

	// 发送第二阶段响应：扫描已结束，进度100%
	nc.cmdResponseChan <- &node_rpc.CmdResponse{
		CmdType: node_rpc.CmdType_cmdNetScanType,
		TaskID:  request.TaskID,
		NodeID:  nc.config.System.Uid,
		NetScanOutMessage: &node_rpc.NetScanOutMessage{
			End:      true, // 已结束
			Progress: 100,  // 进度100%
		},
	}
}
