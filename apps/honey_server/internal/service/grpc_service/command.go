package grpc_service

// File: node_command_service.go
// Description: 实现节点与服务端的grpc双向流命令交互服务，负责管理多节点连接、命令发送与接收、连接生命周期及资源释放，确保并发安全

import (
	"errors"
	"io"
	"sync"
	"time"

	"honey_server/internal/rpc/node_rpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// Command 节点命令交互结构体
// 管理单个节点的grpc流连接、命令请求/响应通道及连接状态
type Command struct {
	ReqChan  chan *node_rpc.CmdRequest          // 命令请求通道，用于接收发送给节点的命令
	ResChan  chan *node_rpc.CmdResponse         // 命令响应通道，用于缓存节点返回的响应
	Server   node_rpc.NodeService_CommandServer // 节点对应的grpc流服务实例
	NodeID   string                             // 节点唯一标识
	stopChan chan struct{}                      // 停止信号通道，用于通知协程退出
	wg       sync.WaitGroup                     // 等待组，用于协调发送/接收协程的退出
	mu       sync.RWMutex                       // 互斥锁，用于保护closed状态的并发访问
	closed   bool                               // 连接是否已关闭的标志
}

var (
	NodeCommandMap = make(map[string]*Command) // 节点命令映射表，键为节点ID，值为对应的Command实例
	mapMutex       sync.RWMutex                // 映射表的读写锁，确保多节点并发操作安全
)

// Command 实现grpc的NodeService_CommandServer接口，处理节点的双向流连接
// 功能：从元数据提取节点ID，创建Command实例并注册到映射表，启动发送/接收协程，监听连接关闭并清理资源
func (s NodeService) Command(stream node_rpc.NodeService_CommandServer) error {
	// 从上下文获取元数据（包含节点标识）
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("缺少元数据(metadata)")
	}

	// 从元数据中提取节点ID
	nodeIDList := md.Get("nodeID")
	if len(nodeIDList) == 0 {
		return errors.New("元数据中未找到nodeID")
	}
	nodeID := nodeIDList[0]

	// 初始化当前节点的Command实例，请求/响应通道设置缓冲避免阻塞
	cmd := &Command{
		ReqChan:  make(chan *node_rpc.CmdRequest, 10),
		ResChan:  make(chan *node_rpc.CmdResponse, 10),
		Server:   stream,
		NodeID:   nodeID,
		stopChan: make(chan struct{}),
	}

	// 加锁保护映射表，添加当前节点的Command实例
	mapMutex.Lock()
	NodeCommandMap[nodeID] = cmd
	mapMutex.Unlock()

	logrus.Infof("节点 %s 已连接", nodeID)

	// 启动发送和接收协程，等待组计数+2
	cmd.wg.Add(2)
	go cmd.sendLoop()    // 负责向节点发送命令
	go cmd.receiveLoop() // 负责接收节点的响应

	// 监听上下文取消（如节点断开连接），触发资源清理
	go func() {
		<-ctx.Done()
		logrus.Infof("节点 %s 的上下文已取消", nodeID)
		cmd.Close()
	}()

	// 等待发送/接收协程完成
	cmd.wg.Wait()

	// 从映射表中移除当前节点的Command实例
	mapMutex.Lock()
	delete(NodeCommandMap, nodeID)
	mapMutex.Unlock()

	logrus.Infof("节点 %s 已断开连接", nodeID)
	return nil
}

// sendLoop 命令发送循环协程
// 功能：从ReqChan读取命令并通过grpc流发送给节点，处理发送错误及停止信号
func (c *Command) sendLoop() {
	defer c.wg.Done() // 协程退出时通知等待组

	for {
		select {
		case req, ok := <-c.ReqChan:
			// 若请求通道已关闭，退出循环
			if !ok {
				logrus.Infof("节点 %s 的请求通道已关闭", c.NodeID)
				return
			}

			// 发送命令到节点
			err := c.Server.Send(req)
			if err != nil {
				logrus.Errorf("向节点 %s 发送命令失败: %v", c.NodeID, err)
				c.Close() // 发送失败时关闭连接
				return
			}

		case <-c.stopChan:
			// 收到停止信号，退出循环
			logrus.Infof("停止节点 %s 的发送协程", c.NodeID)
			return
		}
	}
}

// receiveLoop 响应接收循环协程
// 功能：从grpc流接收节点的响应并写入ResChan，处理接收错误、超时及停止信号
func (c *Command) receiveLoop() {
	defer c.wg.Done() // 协程退出时通知等待组

	for {
		// 从节点接收响应
		res, err := c.Server.Recv()
		if err != nil {
			if err == io.EOF {
				logrus.Infof("节点 %s 正常断开连接", c.NodeID)
			} else {
				logrus.Errorf("从节点 %s 接收响应失败: %v", c.NodeID, err)
			}
			c.Close() // 接收失败时关闭连接
			return
		}

		// 带超时发送响应到ResChan，避免通道阻塞导致协程卡死
		select {
		case c.ResChan <- res:
			logrus.Debugf("节点 %s 的响应已写入通道", c.NodeID)
		case <-time.After(5 * time.Second):
			logrus.Warnf("节点 %s 的响应发送超时，已丢弃", c.NodeID)
		case <-c.stopChan:
			// 收到停止信号，退出循环
			logrus.Infof("停止节点 %s 的接收协程", c.NodeID)
			return
		}
	}
}

// Close 关闭节点的命令交互资源
// 功能：安全关闭停止信号通道和请求通道，标记连接为关闭状态，清空响应通道避免资源泄漏
func (c *Command) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 若已关闭则直接返回
	if c.closed {
		return
	}

	logrus.Infof("关闭节点 %s 的命令通道", c.NodeID)
	c.closed = true // 标记为已关闭

	close(c.stopChan) // 通知发送/接收协程退出
	close(c.ReqChan)  // 关闭请求通道

	// 启动协程清空响应通道，避免阻塞
	go func() {
		for range c.ResChan {
		}
	}()
}

// GetNodeCommand 安全获取节点对应的Command实例
// 功能：通过节点ID从映射表中查询Command实例，使用读锁保证并发安全
// 返回值：节点的Command实例及是否存在的标志
func GetNodeCommand(nodeID string) (*Command, bool) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	cmd, ok := NodeCommandMap[nodeID]
	return cmd, ok
}
