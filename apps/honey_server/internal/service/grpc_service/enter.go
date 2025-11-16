package grpc_service

// File: service/grpc_service/enter.go
// Description: gRPC服务，负责启动gRPC服务器，注册节点服务处理器，处理节点注册等gRPC请求。

import (
	"context"
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/rpc/node_rpc"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// NodeService 节点服务的gRPC处理器结构体
// 实现了node_rpc.NodeServiceServer接口，用于处理节点相关的gRPC请求
type NodeService struct {
	node_rpc.UnimplementedNodeServiceServer
}

// Register 处理节点注册请求的gRPC方法
func (NodeService) Register(ctx context.Context, request *node_rpc.RegisterRequest) (pd *node_rpc.BaseResponse, err error) {
	pd = new(node_rpc.BaseResponse)
	fmt.Println("节点注册", request)
	return
}

// Run 启动gRPC服务的入口函数
//
// 1. 监听配置文件中指定的gRPC端口
// 2. 创建gRPC服务器实例
// 3. 注册节点服务处理器
// 4. 启动服务器并开始处理请求
func Run() {
	// 从全局配置中获取gRPC服务监听地址
	addr := global.Config.System.GrpcAddr
	// 监听指定的TCP地址，准备接收gRPC连接
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatalf("Failed to listen: %v", err)
	}

	// 创建一个新的gRPC服务器实例
	s := grpc.NewServer()
	// 实例化节点服务处理器
	server := NodeService{}
	// 将节点服务处理器注册到gRPC服务器，使其能处理对应的gRPC请求
	node_rpc.RegisterNodeServiceServer(s, &server)
	// 打印服务启动日志
	logrus.Infof("grpc server running %s", addr)
	// 启动gRPC服务器，开始监听并处理客户端请求
	err = s.Serve(listen)
	if err != nil {
		logrus.Fatalf("Failed to serve: %v", err)
	}
}
