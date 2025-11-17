package command

// File: service/command/command_network_list.go
// Description: 节点网卡列表获取及格式转换逻辑，负责根据过滤条件获取网卡信息并转换为RPC通信所需格式

import (
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/info"
)

// getNetworkList 获取节点的网卡列表信息
func (nc *NodeClient) getNetworkList(filters []string) ([]*node_rpc.NetworkInfoMessage, error) {
	// 调用工具函数获取原始网卡列表（应用过滤条件）
	_networkList, err := info.GetNetworkList(filters)
	if err != nil {
		return nil, err
	}

	// 将原始网卡信息转换为RPC通信使用的NetworkInfoMessage格式
	var networkList []*node_rpc.NetworkInfoMessage
	for _, networkInfo := range _networkList {
		networkList = append(networkList, &node_rpc.NetworkInfoMessage{
			Network: networkInfo.Network,     // 网卡接口名称
			Ip:      networkInfo.Ip,          // 网卡接口的IP地址
			Net:     networkInfo.Net,         // 子网信息
			Mask:    int32(networkInfo.Mask), // 子网掩码
		})
	}

	return networkList, nil
}
