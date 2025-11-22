package main

// File: main.go
// Description: 节点程序主入口

import (
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/service/command"
	"honey_node/internal/service/cron_service"
	"honey_node/internal/service/mq_service"

	"github.com/sirupsen/logrus"
)

// nodeClient 全局节点客户端实例
var nodeClient *command.NodeClient

func main() {
	// 初始化系统配置：从配置文件读取全局配置
	global.Config = core.ReadConfig()
	// 初始化日志系统：设置默认日志格式、输出方式等
	core.SetLogDefault()
	// 获取全局日志实例：供全系统使用统一的日志接口
	global.Log = core.GetLogger()

	// 创建gRPC客户端：建立与服务端的gRPC连接，用于后续通信
	global.GrpcClient = core.GetGrpcClient()

	// 初始化节点客户端：封装gRPC通信逻辑，提供节点注册、命令处理等能力
	nodeClient = command.NewNodeClient(global.GrpcClient, global.Config)

	// 节点注册：向服务端注册自身信息，完成节点上线流程
	if err := nodeClient.Register(); err != nil {
		logrus.Fatalf("节点注册失败: %v", err)
		return
	}

	// 初始化消息队列：建立与RabbitMQ的连接，用于消费服务端下发的任务消息
	global.Queue = core.InitMQ()

	// 启动命令处理服务：监听并处理服务端通过gRPC下发的命令
	nodeClient.StartCommandHandling()

	// 启动定时任务服务：执行节点本地的周期性任务
	cron_service.Run()
	// 启动消息队列消费服务：消费RabbitMQ中的任务消息
	mq_service.Run()

	// 阻塞主线程，保持程序运行（避免main函数退出）
	select {}
}
