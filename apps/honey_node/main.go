package main

// File: main.go
// Description: grpc客户端主程序

import (
	"context"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/service/cron_service"
	"honey_node/internal/utils/info"
	"honey_node/internal/utils/ip"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// 读取系统配置
	global.Config = core.ReadConfig()
	// 设置日志默认配置
	core.SetLogDefault()
	// 获取日志实例
	global.Log = core.GetLogger()
	// 初始化grpc客户端
	global.GrpcClient = core.GetGrpcClient()

	// 执行节点注册
	err := register()
	if err != nil {
		logrus.Errorf("节点注册失败 %s", err)
		return
	}
	logrus.Infof("节点注册成功")

	// 启动定时任务调度器
	cron_service.Run()

	// 阻塞当前goroutine，防止程序退出
	select {}
}

// 节点注册函数
func register() (err error) {
	// 获取节点网络信息（IP和MAC地址），基于配置的网络接口
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
	var networkList []*node_rpc.NetworkInfoMessage
	_networkList, err := info.GetNetworkList("hy-")
	if err != nil {
		return
	}
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
		Ip:      _ip,
		Mac:     mac,
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
		NetworkList: networkList,
	}

	// 通过grpc客户端发送注册请求
	_, err = global.GrpcClient.Register(context.Background(), &req)
	if err != nil {
		logrus.Fatalf("节点注册失败 %s", err)
		return
	}
	logrus.Infof("节点注册 上报信息 %v", req)
	return nil
}
