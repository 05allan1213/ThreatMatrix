package command

// File: service/command/enter.go
// Description: 节点客户端核心逻辑实现，负责管理与服务器的grpc连接、命令收发、连接维护及重连机制

import (
	"context"
	"honey_node/internal/config"
	"honey_node/internal/rpc/node_rpc"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NodeClient 节点客户端结构体，用于管理与服务器的grpc连接及命令处理
type NodeClient struct {
	client          node_rpc.NodeServiceClient         // grpc客户端实例，用于与服务器通信
	config          *config.Config                     // 节点配置信息
	cmdResponseChan chan *node_rpc.CmdResponse         // 命令响应通道，用于向服务器发送处理结果
	stream          node_rpc.NodeService_CommandClient // 命令流，用于双向通信
	ctx             context.Context                    // 上下文，用于控制协程生命周期
	cancel          context.CancelFunc                 // 上下文取消函数，用于终止相关协程
	wg              sync.WaitGroup                     // 等待组，用于协调多个协程的退出
	reconnectTimer  *time.Timer                        // 重连计时器，用于连接断开后的延迟重连
	mu              sync.Mutex                         // 互斥锁，用于保护共享资源（如连接状态）
	isConnected     bool                               // 连接状态标识，true表示已连接
}

// NewNodeClient 创建NodeClient实例
func NewNodeClient(grpcClient node_rpc.NodeServiceClient,
	config *config.Config) *NodeClient {
	return &NodeClient{
		client:          grpcClient,
		config:          config,
		cmdResponseChan: make(chan *node_rpc.CmdResponse, 10), // 带缓冲的响应通道，避免阻塞
		reconnectTimer:  time.NewTimer(0),
	}
}

// StartCommandHandling 启动命令处理机制和自动重连逻辑
// 初始化上下文，启动连接协程，负责维护与服务器的连接
func (nc *NodeClient) StartCommandHandling() {
	nc.ctx, nc.cancel = context.WithCancel(context.Background())
	nc.wg.Add(1)

	go func() {
		defer nc.wg.Done()

		// 初始连接服务器
		nc.connect()

		// 监听上下文取消信号，退出时断开连接
		<-nc.ctx.Done()
		nc.disconnect()
		logrus.Info("命令处理已停止")
	}()
}

// connect 建立与服务器的命令流连接
// 加锁保证线程安全，创建grpc流，启动收发协程
func (nc *NodeClient) connect() {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	// 若已连接则直接返回
	if nc.isConnected {
		return
	}

	// 构建包含节点唯一标识的元数据（用于服务器身份验证）
	ctx := metadata.NewOutgoingContext(nc.ctx, metadata.Pairs("nodeID", nc.config.System.Uid))

	// 创建与服务器的双向命令流
	stream, err := nc.client.Command(ctx)
	if err != nil {
		logrus.Errorf("创建命令流失败: %v，将在2秒后重试", err)
		nc.scheduleReconnect(2 * time.Second) // 连接失败，安排重连
		return
	}

	nc.stream = stream
	nc.isConnected = true
	logrus.Info("节点命令流连接成功")

	// 启动响应发送和请求接收协程
	nc.wg.Add(2)
	go nc.sendResponses()
	go nc.receiveRequests()
}

// disconnect 断开与服务器的连接
// 释放资源，重置连接状态，停止重连计时器
func (nc *NodeClient) disconnect() {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	if !nc.isConnected {
		return
	}

	// 停止重连计时器
	nc.reconnectTimer.Stop()

	// 关闭命令流的发送端
	if nc.stream != nil {
		nc.stream.CloseSend()
		nc.stream = nil
	}

	// 关闭并重建响应通道（避免通道残留数据影响下次连接）
	close(nc.cmdResponseChan)
	nc.cmdResponseChan = make(chan *node_rpc.CmdResponse, 10)

	nc.isConnected = false
	logrus.Info("节点命令流已断开")
}

// scheduleReconnect 安排延迟重连
func (nc *NodeClient) scheduleReconnect(delay time.Duration) {
	nc.reconnectTimer.Reset(delay)

	go func() {
		<-nc.reconnectTimer.C
		// 若上下文已取消（如节点退出），则不再重连
		if nc.ctx.Err() != nil {
			return
		}
		nc.connect() // 尝试重新连接
	}()
}

// sendResponses 向服务器发送命令响应
// 循环从响应通道读取数据并发送，处理发送失败时的重连逻辑
func (nc *NodeClient) sendResponses() {
	defer nc.wg.Done()

	for {
		select {
		case <-nc.ctx.Done():
			// 上下文取消，退出协程
			return

		case response, ok := <-nc.cmdResponseChan:
			if !ok {
				// 通道已关闭，退出协程
				return
			}

			// 发送响应到服务器
			if err := nc.stream.Send(response); err != nil {
				logrus.Errorf("发送响应失败: %v", err)
				nc.disconnect()                       // 发送失败，断开连接
				nc.scheduleReconnect(2 * time.Second) // 安排重连
				return
			}

			logrus.Debugf("已发送响应: %+v", response)
		}
	}
}

// receiveRequests 从服务器接收命令并处理
// 循环接收命令，分发到对应的处理函数，处理接收失败时的重连逻辑
func (nc *NodeClient) receiveRequests() {
	defer nc.wg.Done()

	for {
		// 从命令流接收服务器发送的命令
		request, err := nc.stream.Recv()
		if err != nil {
			// 解析错误类型，输出对应日志
			if status.Code(err) == 0 { // io.EOF 错误（服务器主动关闭连接）
				logrus.Info("服务器关闭了连接")
			} else if err == context.Canceled {
				logrus.Info("上下文已取消")
			} else if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				logrus.Warnf("临时网络错误: %v", err)
			} else {
				logrus.Errorf("接收请求失败: %v", err)
			}

			nc.disconnect()                       // 接收失败，断开连接
			nc.scheduleReconnect(2 * time.Second) // 安排重连
			return
		}

		logrus.Infof("收到命令: %+v", request)
		nc.handleCommand(request) // 处理接收到的命令
	}
}

// handleCommand 处理服务器发送的命令
// 根据命令类型分发到对应的处理方法
func (nc *NodeClient) handleCommand(request *node_rpc.CmdRequest) {
	switch request.CmdType {
	case node_rpc.CmdType_cmdNetworkFlushType:
		// 处理网络刷新命令
		nc.CmdNetworkFlush(request)
	case node_rpc.CmdType_cmdNetScanType:
		// 处理网络扫描命令
		nc.CmdNetScan(request)
	default:
		// 未知命令类型，记录警告日志
		logrus.Warnf("未知命令类型: %v", request.CmdType)
	}
}
