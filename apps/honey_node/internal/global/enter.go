package global

// File: global/enter.go
// Description: 声明全局变量，供其他模块共享。

import (
	"honey_node/internal/config"
	"honey_node/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	Version   = "v1.0.1"              // 版本号
	Commit    = "56a9c63f"            // 提交ID
	BuildTime = "2025-11-10 17:12:34" // 构建时间
)

var (
	Config     *config.Config             // 配置实例
	Log        *logrus.Entry              // 日志实例
	GrpcClient node_rpc.NodeServiceClient // RPC客户端实例
	Queue      *amqp.Channel              // rabbitMQ消息队列实例
)
