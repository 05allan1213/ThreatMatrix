package core

// File: core/grpc_client.go
// Description: 提供获取节点服务grpc客户端实例的功能

import (
	"honey_node/internal/global"
	"honey_node/internal/rpc"
	"honey_node/internal/rpc/node_rpc"
)

// 获取节点服务的grpc客户端实例
func GetGrpcClient() node_rpc.NodeServiceClient {
	// 从全局配置中获取grpc管理服务的地址
	addr := global.Config.System.GrpcManageAddr

	// 调用rpc包的GetConn方法创建与grpc服务端的连接
	conn := rpc.GetConn(addr)

	// 初始化节点服务grpc客户端
	client := node_rpc.NewNodeServiceClient(conn)

	return client
}
