package main

// File: main.go
// Description: grpc节点程序主入口

import (
	"context"
	"fmt"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/service/cron_service"
	"honey_node/internal/utils/info"
	"honey_node/internal/utils/ip"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

func main() {
	// 读取系统配置
	global.Config = core.ReadConfig()
	// 设置日志默认配置
	core.SetLogDefault()
	// 获取日志实例
	global.Log = core.GetLogger()
	// 初始化grpc客户端（连接管理服务）
	global.GrpcClient = core.GetGrpcClient()

	// 执行节点注册流程
	err := register()
	if err != nil {
		logrus.Errorf("节点注册失败 %s", err)
		return
	}
	logrus.Infof("节点注册成功")

	// 启动命令交互协程（处理与管理服务的双向流通信）
	go command()

	// 启动定时任务调度器
	cron_service.Run()

	// 阻塞当前goroutine，保持程序运行
	select {}
}

// CmdResponseChan 命令响应通道
// 用于缓存需要发送给管理服务的命令执行结果，由命令处理逻辑写入，发送协程读取并通过grpc流发送
var CmdResponseChan = make(chan *node_rpc.CmdResponse)

// Stream grpc命令流客户端实例
// 用于与管理服务建立双向流连接，发送命令响应和接收管理服务的命令
var Stream node_rpc.NodeService_CommandClient

// command 处理与管理服务的grpc双向流命令交互
// 功能：建立命令流连接，启动响应发送协程，循环接收管理服务命令并处理，处理结果通过响应通道返回；连接断开时自动重连
func command() {
	// 创建包含节点唯一标识（UID）的上下文（用于grpc身份验证）
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("nodeID", global.Config.System.Uid))

	// 建立与管理服务的命令流连接
	var err error
	Stream, err = global.GrpcClient.Command(ctx)
	if err != nil {
		logrus.Errorf("节点Command流连接失败 %s", err)
		// 连接失败后等待2秒重试
		time.Sleep(2 * time.Second)
		command()
		return
	}

	// 启动协程：从响应通道读取结果并通过grpc流发送给管理服务
	go func() {
		for response := range CmdResponseChan {
			err := Stream.Send(response)
			if err != nil {
				logrus.Infof("命令响应发送失败 %s", err)
				continue
			}
		}
	}()

	fmt.Println("节点成功连接到管理服务的Command流")

	// 循环接收管理服务发送的命令
	for {
		request, err := Stream.Recv()
		if err == io.EOF {
			// 流结束（管理服务断开连接）
			logrus.Infof("与管理服务的Command流断开（EOF）")
			break
		}
		if err != nil {
			// 接收命令出错
			logrus.Infof("接收管理服务命令出错 %s", err)
			break
		}

		// 打印接收的命令
		fmt.Println("接收管理服务的命令数据", request)

		// 根据命令类型处理不同逻辑
		switch request.CmdType {
		case node_rpc.CmdType_cmdNetworkFlushType:
			// 处理网卡信息刷新命令
			fmt.Println("执行网卡信息刷新命令")

			// 获取过滤后的网卡列表（过滤指定前缀的网卡）
			_networkList, _ := info.GetNetworkList(request.NetworkFlushInMessage.FilterNetworkName[0])

			// 转换网卡信息为grpc消息格式
			var networkList []*node_rpc.NetworkInfoMessage
			for _, networkInfo := range _networkList {
				networkList = append(networkList, &node_rpc.NetworkInfoMessage{
					Network: networkInfo.Network,     // 网卡名称
					Ip:      networkInfo.Ip,          // 网卡IP地址
					Net:     networkInfo.Net,         // 网络地址（CIDR格式）
					Mask:    int32(networkInfo.Mask), // 子网掩码长度
				})
			}

			// 构建命令响应并写入响应通道（使用原请求的TaskID保持关联）
			CmdResponseChan <- &node_rpc.CmdResponse{
				CmdType: node_rpc.CmdType_cmdNetworkFlushType, // 响应命令类型（与请求一致）
				TaskID:  request.TaskID,                       // 任务ID（与请求一致，用于关联命令和结果）
				NodeID:  global.Config.System.Uid,             // 节点唯一标识
				NetworkFlushOutMessage: &node_rpc.NetworkFlushOutMessage{
					NetworkList: networkList, // 刷新后的网卡列表
				},
			}
		}
	}

	// 连接断开后，等待2秒自动重连
	time.Sleep(2 * time.Second)
	command()
}

// register 节点注册函数
// 功能：收集节点网络信息（IP、MAC）、主机名、系统信息、网卡列表，构建注册请求并发送给管理服务，完成节点注册
func register() (err error) {
	// 获取节点指定网络接口的IP和MAC地址
	_ip, mac, err := ip.GetNetworkInfo(global.Config.System.Network)
	if err != nil {
		return
	}

	// 获取节点主机名
	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	// 获取节点系统信息（操作系统版本、内核版本等）
	systemInfo, err := info.GetSystemInfo()
	if err != nil {
		return
	}

	// 获取节点网卡列表（过滤"hy-"前缀的网卡）
	var networkList []*node_rpc.NetworkInfoMessage
	_networkList, err := info.GetNetworkList("hy-")
	if err != nil {
		return
	}
	// 转换网卡列表为grpc消息格式
	for _, networkInfo := range _networkList {
		networkList = append(networkList, &node_rpc.NetworkInfoMessage{
			Network: networkInfo.Network,
			Ip:      networkInfo.Ip,
			Net:     networkInfo.Net,
			Mask:    int32(networkInfo.Mask),
		})
	}

	// 构建节点注册请求
	req := node_rpc.RegisterRequest{
		Ip:      _ip,                      // 节点IP地址
		Mac:     mac,                      // 节点MAC地址
		NodeUid: global.Config.System.Uid, // 节点唯一标识
		Version: global.Version,           // 节点程序版本
		Commit:  global.Commit,            // 节点提交哈希
		SystemInfo: &node_rpc.SystemInfoMessage{
			HostName:            hostname,                // 主机名
			DistributionVersion: systemInfo.OSVersion,    // 操作系统版本
			CoreVersion:         systemInfo.Kernel,       // 内核版本
			SystemType:          systemInfo.Architecture, // 系统架构
			StartTime:           systemInfo.BootTime,     // 系统启动时间
		},
		NetworkList: networkList, // 节点网卡列表
	}

	// 发送注册请求到管理服务
	_, err = global.GrpcClient.Register(context.Background(), &req)
	if err != nil {
		logrus.Fatalf("节点注册失败 %s", err)
		return
	}
	logrus.Infof("节点注册 上报信息 %v", req)
	return nil
}
