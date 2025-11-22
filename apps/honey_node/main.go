package main

// File: main.go
// Description: 节点程序主入口

import (
	"context"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/service/command"
	"honey_node/internal/service/cron_service"
	"honey_node/internal/service/mq_service"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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

	// 启动定时任务服务：执行节点本地的周期性任务（如心跳上报、状态检测等）
	cron_service.Run()
	// 启动消息队列消费服务：消费RabbitMQ中的任务消息（如创建IP、删除IP等）
	mq_service.Run()

	// 启动TCP监听服务（异步goroutine）：建立本地TCP端口监听，通过gRPC隧道转发数据
	go tcpListen()

	// 阻塞主线程，保持程序运行（避免main函数退出）
	select {}
}

// 本地监听地址与目标转发地址配置
// localAddr：节点本地监听的TCP地址，接收外部请求
// targetAddr：通过gRPC隧道转发到的目标地址（服务端侧）
var localAddr = "0.0.0.0:8005"
var targetAddr = "127.0.0.1:8000"

// tcpListen 启动本地TCP监听服务
// 接收本地客户端连接，并通过gRPC隧道转发到服务端目标地址，实现端口映射/数据透传
func tcpListen() {
	// 创建本地TCP监听
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("创建本地监听失败: %v", err)
	}
	defer listener.Close() // 函数退出时关闭监听，释放端口资源

	log.Printf("本地监听启动，地址: %s", localAddr)
	log.Printf("目标地址: %s", targetAddr)

	// 信号处理：监听系统中断信号（Ctrl+C、kill），实现优雅关闭
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		log.Println("接收到终止信号，优雅关闭...")
		os.Exit(0)
	}()

	// 循环接受客户端连接：为每个连接创建独立goroutine处理，支持并发连接
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("接受客户端连接失败: %v", err)
			continue
		}

		// 异步处理单个连接的转发逻辑，避免阻塞监听主循环
		go handleConnection(global.GrpcClient, clientConn, targetAddr)
	}
}

// handleConnection 处理单个TCP连接的转发逻辑
// 通过gRPC Tunnel双向流将本地TCP连接的数据转发到服务端目标地址，实现隧道透传
func handleConnection(client node_rpc.NodeServiceClient, localConn net.Conn, targetAddr string) {
	defer localConn.Close() // 连接处理完成后关闭本地TCP连接

	// 创建带取消功能的上下文：用于控制gRPC流的生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建gRPC Tunnel双向流：建立与服务端的隧道连接
	stream, err := client.Tunnel(ctx)
	if err != nil {
		log.Printf("创建隧道失败: %v", err)
		return
	}

	// 发送初始消息：告知服务端目标转发地址（触发服务端连接目标地址）
	if err := stream.Send(&node_rpc.TunnelData{
		Chunk:   []byte{}, // 初始消息无数据，仅传递地址
		Address: targetAddr,
	}); err != nil {
		log.Printf("发送初始请求失败: %v", err)
		return
	}

	// 启动goroutine处理"服务端→gRPC流→本地TCP连接"的数据流向
	// 从gRPC流读取服务端转发的数据，写入本地客户端连接
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF { // 服务端关闭流则退出
				break
			}
			if err != nil {
				log.Printf("接收gRPC服务器数据失败: %v", err)
				break
			}

			// 将服务端转发的数据写入本地TCP连接
			_, err = localConn.Write(resp.Chunk)
			if err != nil {
				log.Printf("写入本地连接失败: %v", err)
				break
			}
		}
		cancel() // 数据流向中断时取消上下文，终止另一方向的处理
	}()

	// 处理"本地TCP连接→gRPC流→服务端"的数据流向
	// 从本地客户端连接读取数据，通过gRPC流发送到服务端
	buffer := make([]byte, 4096) // 4KB缓冲区平衡IO性能与内存占用
	for {
		n, err := localConn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("本地连接已关闭")
			} else {
				log.Printf("从本地连接读取失败: %v", err)
			}
			break
		}

		// 将本地连接的数据通过gRPC流发送到服务端
		err = stream.Send(&node_rpc.TunnelData{
			Chunk:   buffer[:n], // 仅发送实际读取的字节
			Address: targetAddr,
		})
		if err != nil {
			log.Printf("发送数据到gRPC服务器失败: %v", err)
			break
		}
	}

	// 主动关闭gRPC流的发送端，告知服务端数据传输完成
	stream.CloseSend()
}
