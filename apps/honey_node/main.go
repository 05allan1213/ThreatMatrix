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
	// 读取系统配置
	global.Config = core.ReadConfig()
	// 设置日志默认配置
	core.SetLogDefault()
	// 获取日志实例
	global.Log = core.GetLogger()

	// 创建grpc客户端
	global.GrpcClient = core.GetGrpcClient()

	// 初始化节点客户端
	nodeClient = command.NewNodeClient(global.GrpcClient, global.Config)

	// 执行节点注册
	if err := nodeClient.Register(); err != nil {
		logrus.Fatalf("节点注册失败: %v", err)
		return
	}

	// 初始化rabbitMQ消息队列
	global.Queue = core.InitMQ()

	// 启动命令处理机制（建立与服务器的双向流连接，处理命令收发及自动重连）
	nodeClient.StartCommandHandling()

	// 启动定时任务
	cron_service.Run()

	// 启动消息队列服务
	mq_service.Run()

	// 阻塞主goroutine（防止程序退出，保持所有后台协程运行）
	select {}
}
