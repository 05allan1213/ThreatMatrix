package grpc_service

// File: command_stream.go
// Description: 实现grpc双向流命令服务，负责处理节点与管理服务之间的命令交互，通过通道实现命令的发送与响应接收

import (
	"errors"
	"fmt"
	"honey_server/internal/rpc/node_rpc"
	"io"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

// CmdRequestChan 命令请求通道
// 用于传递需要发送给节点的命令请求，由API层写入，grpc流服务读取并发送到节点
var CmdRequestChan = make(chan *node_rpc.CmdRequest)

// CmdResponseChan 命令响应通道
// 用于传递节点返回的命令执行结果，由grpc流服务接收并写入，API层读取并返回给前端
var CmdResponseChan = make(chan *node_rpc.CmdResponse)

// StreamMap 节点流映射表
var StreamMap = map[string]node_rpc.NodeService_CommandServer{}

// Command 实现grpc的双向流命令接口
func (NodeService) Command(stream node_rpc.NodeService_CommandServer) error {
	ctx := stream.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	nodeIDList := md.Get("nodeID")
	if len(nodeIDList) == 0 {
		return errors.New("请在metadata中传入节点id")
	}
	nodeID := nodeIDList[0]
	StreamMap[nodeID] = stream

	// 启动goroutine接收节点发送的命令响应
	go func() {
		for {
			// 从流中接收节点返回的响应
			response, err := StreamMap[nodeID].Recv()
			if err == io.EOF {
				// 流结束（节点断开连接），退出接收循环
				logrus.Infof("节点断开")
				return
			}
			if err != nil {
				// 接收出错，记录日志并退出接收循环
				logrus.Infof("节点出错 %s", err)
				return
			}
			// 打印命令执行结果（调试用）
			fmt.Println("命令结果", response)
			// 将响应发送到响应通道，供API层读取
			CmdResponseChan <- response
		}
	}()

	// 从命令请求通道读取命令，并发送到节点
	for request := range CmdRequestChan {
		err := StreamMap[nodeID].Send(request)
		if err != nil {
			// 发送失败，记录日志并继续处理下一个请求
			logrus.Infof("数据发送失败 %s", err)
			continue
		}
	}
	return nil
}
