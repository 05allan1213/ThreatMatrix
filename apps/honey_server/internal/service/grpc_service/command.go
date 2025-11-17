package grpc_service

// File: service/grpc_service/command.go
// Description: 管理节点与服务端的grpc双向流命令交互，通过节点ID映射命令通道，实现多节点的命令发送与响应接收

import (
	"errors"
	"fmt"
	"honey_server/internal/rpc/node_rpc"
	"io"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// Command 节点命令交互结构体
// 用于管理单个节点的命令请求、响应通道及grpc流连接
type Command struct {
	ReqChan chan *node_rpc.CmdRequest          // 命令请求通道，用于接收需要发送给节点的命令
	ResChan chan *node_rpc.CmdResponse         // 命令响应通道，用于缓存节点返回的命令执行结果
	Server  node_rpc.NodeService_CommandServer // 节点对应的grpc流服务实例
}

// NodeCommandMap 节点命令映射表
// 键为节点ID，值为对应节点的Command实例，用于管理多个节点的命令流交互
var NodeCommandMap = map[string]*Command{}

// Command 实现grpc的双向流命令接口，支持多节点并发命令交互
// 功能：通过metadata获取节点ID，创建并存储节点对应的命令通道，启动协程发送命令，
func (NodeService) Command(stream node_rpc.NodeService_CommandServer) error {
	// 从上下文获取元数据（包含节点ID）
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil // 元数据获取失败，直接返回
	}

	// 从元数据中提取节点ID
	nodeIDList := md.Get("nodeID")
	if len(nodeIDList) == 0 {
		return errors.New("请在metadata中传入节点id") // 节点ID不存在，返回错误
	}
	nodeID := nodeIDList[0]

	// 为当前节点创建Command实例并存储到映射表
	NodeCommandMap[nodeID] = &Command{
		ReqChan: make(chan *node_rpc.CmdRequest),
		ResChan: make(chan *node_rpc.CmdResponse),
		Server:  stream,
	}

	// 启动协程：从请求通道读取命令并发送到节点
	go func() {
		for request := range NodeCommandMap[nodeID].ReqChan {
			err := NodeCommandMap[nodeID].Server.Send(request)
			if err != nil {
				logrus.Infof("向节点[%s]发送命令失败 %s", nodeID, err)
				continue
			}
		}
	}()

	// 循环接收节点发送的命令响应
	for {
		response, err := NodeCommandMap[nodeID].Server.Recv()
		if err == io.EOF {
			// 流结束（节点断开连接）
			logrus.Infof("节点[%s]断开连接", nodeID)
			break
		}
		if err != nil {
			// 接收响应出错
			logrus.Infof("接收节点[%s]响应出错 %s", nodeID, err)
			break
		}

		// 打印命令执行结果（调试用）
		fmt.Printf("节点[%s]命令结果: %v\n", nodeID, response)

		// 将响应发送到当前节点的响应通道
		NodeCommandMap[nodeID].ResChan <- response
	}

	// 节点断开后，从映射表中删除该节点的Command实例
	delete(NodeCommandMap, nodeID)
	return nil
}
