package command

// File: server/command/enter.go
// Description: 负责节点与服务器的双向命令交互

import (
	"context"
	"fmt"
	"honey_node/internal/config"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/info"
	"honey_node/internal/utils/ip"
	"net"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NodeClient 节点客户端结构体
type NodeClient struct {
	client          node_rpc.NodeServiceClient         // grpc客户端实例，用于与服务器通信
	config          *config.Config                     // 节点配置信息
	cmdResponseChan chan *node_rpc.CmdResponse         // 命令响应通道，缓存需要发送给服务器的响应
	stream          node_rpc.NodeService_CommandClient // 命令流客户端，用于双向流通信
	ctx             context.Context                    // 上下文，用于控制连接生命周期
	cancel          context.CancelFunc                 // 上下文取消函数，用于终止连接
	wg              sync.WaitGroup                     // 等待组，协调各协程退出
	reconnectTimer  *time.Timer                        // 重连定时器，用于连接断开后的重试
	mu              sync.Mutex                         // 互斥锁，保护连接状态等共享资源的并发访问
	isConnected     bool                               // 连接状态标识，true表示已连接
}

// 创建NodeClient实例
func NewNodeClient(grpcClient node_rpc.NodeServiceClient,
	config *config.Config) *NodeClient {
	return &NodeClient{
		client:          grpcClient,
		config:          config,
		cmdResponseChan: make(chan *node_rpc.CmdResponse, 10), // 带缓冲的响应通道，避免阻塞
		reconnectTimer:  time.NewTimer(0),                     // 初始化重连定时器
	}
}

// 向服务器注册节点
func (nc *NodeClient) Register() error {
	// 创建带10秒超时的上下文，控制注册请求的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取节点指定网络接口的IP和MAC地址
	_ip, mac, err := ip.GetNetworkInfo(nc.config.System.Network)
	if err != nil {
		return fmt.Errorf("获取网络信息失败: %v", err)
	}

	// 获取节点主机名
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("获取主机名失败: %v", err)
	}

	// 获取节点系统信息（操作系统版本、内核等）
	systemInfo, err := info.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("获取系统信息失败: %v", err)
	}

	// 获取节点网卡列表（过滤指定前缀的网卡）
	networkList, err := nc.getNetworkList(nc.config.FilterNetworkList)
	if err != nil {
		return fmt.Errorf("获取网络列表失败: %v", err)
	}

	// 构建节点注册请求
	req := &node_rpc.RegisterRequest{
		Ip:      _ip,
		Mac:     mac,
		NodeUid: nc.config.System.Uid, // 节点唯一标识
		Version: global.Version,       // 节点程序版本
		Commit:  global.Commit,        // 节点提交哈希
		SystemInfo: &node_rpc.SystemInfoMessage{
			HostName:            hostname,                // 主机名
			DistributionVersion: systemInfo.OSVersion,    // 操作系统版本
			CoreVersion:         systemInfo.Kernel,       // 内核版本
			SystemType:          systemInfo.Architecture, // 系统架构
			StartTime:           systemInfo.BootTime,     // 系统启动时间
		},
		NetworkList: networkList, // 节点网卡列表
	}

	// 发送注册请求到服务器
	_, err = nc.client.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("注册请求失败: %v", err)
	}

	logrus.Infof("节点注册成功，上报信息: %+v", req)
	return nil
}

// 启动命令处理和重连机制
func (nc *NodeClient) StartCommandHandling() {
	// 创建可取消的上下文，用于控制所有子协程
	nc.ctx, nc.cancel = context.WithCancel(context.Background())
	nc.wg.Add(1)

	// 启动主连接协程
	go func() {
		defer nc.wg.Done()

		// 初始连接服务器
		nc.connect()

		// 等待上下文取消信号（如程序退出）
		<-nc.ctx.Done()
		nc.disconnect() // 断开连接并清理资源
		logrus.Info("命令处理已停止")
	}()
}

// 建立与服务器的命令流连接
func (nc *NodeClient) connect() {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	// 若已连接则直接返回
	if nc.isConnected {
		return
	}

	// 创建包含节点ID的元数据上下文，用于服务器识别节点
	ctx := metadata.NewOutgoingContext(nc.ctx, metadata.Pairs("nodeID", nc.config.System.Uid))

	// 建立与服务器的命令流连接
	stream, err := nc.client.Command(ctx)
	if err != nil {
		logrus.Errorf("创建命令流失败: %v，将在2秒后重试", err)
		nc.scheduleReconnect(2 * time.Second) // 安排2秒后重连
		return
	}

	// 更新连接状态和流实例
	nc.stream = stream
	nc.isConnected = true
	logrus.Info("节点命令流连接成功")

	// 启动响应发送和命令接收协程
	nc.wg.Add(2)
	go nc.sendResponses()   // 负责发送响应到服务器
	go nc.receiveRequests() // 负责接收服务器的命令
}

// 断开与服务器的连接
func (nc *NodeClient) disconnect() {
	nc.mu.Lock()
	defer nc.mu.Unlock()

	// 若未连接则直接返回
	if !nc.isConnected {
		return
	}

	// 停止重连定时器
	nc.reconnectTimer.Stop()

	// 关闭命令流的发送端
	if nc.stream != nil {
		nc.stream.CloseSend()
		nc.stream = nil
	}

	// 关闭并重建响应通道（避免通道关闭后无法使用）
	close(nc.cmdResponseChan)
	nc.cmdResponseChan = make(chan *node_rpc.CmdResponse, 10)

	// 更新连接状态
	nc.isConnected = false
	logrus.Info("节点命令流已断开")
}

// 尝试连接重连
func (nc *NodeClient) scheduleReconnect(delay time.Duration) {
	nc.reconnectTimer.Reset(delay)

	go func() {
		<-nc.reconnectTimer.C // 等待定时器到期
		// 检查上下文是否已取消（如程序退出）
		if nc.ctx.Err() != nil {
			return
		}
		nc.connect() // 尝试重新连接
	}()
}

// 发送响应到服务器
func (nc *NodeClient) sendResponses() {
	defer nc.wg.Done()

	for {
		select {
		case <-nc.ctx.Done():
			// 上下文已取消，退出协程
			return

		case response, ok := <-nc.cmdResponseChan:
			// 响应通道已关闭，退出协程
			if !ok {
				return
			}

			// 发送响应到服务器
			if err := nc.stream.Send(response); err != nil {
				logrus.Errorf("发送响应失败: %v", err)
				nc.disconnect()                       // 断开连接
				nc.scheduleReconnect(2 * time.Second) // 安排2秒后重连
				return
			}

			logrus.Debugf("已发送响应: %+v", response)
		}
	}
}

// 接收服务器命令
func (nc *NodeClient) receiveRequests() {
	defer nc.wg.Done()

	for {
		// 从命令流接收命令
		request, err := nc.stream.Recv()
		if err != nil {
			// 根据错误类型处理
			if status.Code(err) == 0 { // io.EOF错误（服务器正常关闭连接）
				logrus.Info("服务器关闭了连接")
			} else if err == context.Canceled { // 上下文已取消
				logrus.Info("上下文已取消")
			} else if netErr, ok := err.(net.Error); ok && netErr.Temporary() { // 临时网络错误
				logrus.Warnf("临时网络错误: %v", err)
			} else { // 其他错误
				logrus.Errorf("接收请求失败: %v", err)
			}

			nc.disconnect()                       // 断开连接
			nc.scheduleReconnect(2 * time.Second) // 安排2秒后重连
			return
		}

		logrus.Infof("收到命令: %+v", request)
		nc.handleCommand(request) // 处理收到的命令
	}
}

// 处理服务器命令
func (nc *NodeClient) handleCommand(request *node_rpc.CmdRequest) {
	switch request.CmdType {
	case node_rpc.CmdType_cmdNetworkFlushType:
		logrus.Info("处理网卡刷新命令")

		// 提取命令中的过滤条件
		var filters []string
		if request.NetworkFlushInMessage != nil && len(request.NetworkFlushInMessage.FilterNetworkName) > 0 {
			filters = request.NetworkFlushInMessage.FilterNetworkName
		}

		// 获取过滤后的网卡列表
		networkList, err := nc.getNetworkList(filters)
		if err != nil {
			logrus.Errorf("获取网络列表失败: %v", err)
			return
		}

		// 构建命令响应（使用原请求的TaskID保持关联）
		response := &node_rpc.CmdResponse{
			CmdType: node_rpc.CmdType_cmdNetworkFlushType, // 响应类型与命令类型一致
			TaskID:  request.TaskID,                       // 任务ID，关联命令与响应
			NodeID:  nc.config.System.Uid,                 // 节点唯一标识
			NetworkFlushOutMessage: &node_rpc.NetworkFlushOutMessage{
				NetworkList: networkList, // 刷新后的网卡列表
			},
		}

		// 将响应加入发送队列，带上下文取消检查
		select {
		case nc.cmdResponseChan <- response:
			logrus.Debugf("已将响应加入发送队列: %+v", response)
		case <-nc.ctx.Done():
			logrus.Warn("上下文已取消，丢弃响应")
		}

	default:
		// 处理未知命令类型
		logrus.Warnf("未知命令类型: %v", request.CmdType)
	}
}

// 获取网络列表并转换为grpc消息格式
func (nc *NodeClient) getNetworkList(filters []string) ([]*node_rpc.NetworkInfoMessage, error) {
	// 获取过滤后的网络接口信息
	_networkList, err := info.GetNetworkList(filters)
	if err != nil {
		return nil, err
	}

	// 转换为grpc消息格式的网络列表
	var networkList []*node_rpc.NetworkInfoMessage
	for _, networkInfo := range _networkList {
		networkList = append(networkList, &node_rpc.NetworkInfoMessage{
			Network: networkInfo.Network,     // 网卡名称
			Ip:      networkInfo.Ip,          // IP地址
			Net:     networkInfo.Net,         // 网络地址（CIDR格式）
			Mask:    int32(networkInfo.Mask), // 子网掩码长度
		})
	}

	return networkList, nil
}
