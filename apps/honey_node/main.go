package main

// File: main.go
// Description: gRPC客户端主程序

import (
	"context"
	"fmt"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 读取配置文件到全局配置变量
	global.Config = core.ReadConfig()
	// 设置默认日志配置
	core.SetLogDefault()
	// 获取日志实例
	global.Log = core.GetLogger()

	// 从配置中获取管理节点的gRPC服务地址
	addr := global.Config.System.GrpcManageAddr

	// 创建gRPC客户端连接
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// 连接失败时打印错误并退出程序
		logrus.Fatalf("%s", fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", addr, err))
	}
	// 延迟关闭连接，确保程序退出时释放资源
	defer conn.Close()

	// 初始化节点服务的gRPC客户端实例
	client := node_rpc.NewNodeServiceClient(conn)

	// 发送节点注册请求到gRPC服务器
	result, err := client.Register(context.Background(), &node_rpc.RegisterRequest{
		Ip:      "",    // 节点IP
		Mac:     "xx",  // 节点MAC地址
		NodeUid: "xxx", // 节点唯一标识
		Version: "",    // 节点版本号
		Commit:  "",    // 代码提交哈希
	})

	// 打印注册结果和可能的错误
	fmt.Println(result, err)
}
