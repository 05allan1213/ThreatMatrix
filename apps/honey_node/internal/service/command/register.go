package command

// File: service/command/register.go
// Description: 节点向服务器注册的逻辑实现，负责收集节点网卡信息、系统信息等，构建注册请求并发送至服务器

import (
	"context"
	"fmt"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/info"
	"honey_node/internal/utils/ip"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Register 向服务器注册节点
// 收集节点的网卡信息、系统信息等，构建注册请求并发送至服务器，完成节点注册流程
func (nc *NodeClient) Register() error {
	// 创建带10秒超时的上下文，确保注册请求不会无限阻塞
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取节点的IP和MAC地址（基于配置的网卡接口）
	_ip, mac, err := ip.GetNetworkInfo(nc.config.System.Network)
	if err != nil {
		return fmt.Errorf("获取网卡信息失败: %v", err)
	}

	// 获取当前主机的主机名
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("获取主机名失败: %v", err)
	}

	// 收集系统信息（如操作系统版本、内核版本等）
	systemInfo, err := info.GetSystemInfo()
	if err != nil {
		return fmt.Errorf("获取系统信息失败: %v", err)
	}

	// 获取节点的网卡列表信息（应用配置的过滤条件）
	networkList, err := nc.getNetworkList(nc.config.FilterNetworkList)
	if err != nil {
		return fmt.Errorf("获取网卡列表失败: %v", err)
	}

	// 构建注册请求参数
	req := &node_rpc.RegisterRequest{
		Ip:      _ip,                  // 节点IP地址
		Mac:     mac,                  // 节点MAC地址
		NodeUid: nc.config.System.Uid, // 节点唯一标识
		Version: global.Version,       // 节点程序版本号
		Commit:  global.Commit,        // 节点提交哈希
		SystemInfo: &node_rpc.SystemInfoMessage{ // 系统信息详情
			HostName:            hostname,                // 主机名
			DistributionVersion: systemInfo.OSVersion,    // 操作系统发行版本
			CoreVersion:         systemInfo.Kernel,       // 内核版本
			SystemType:          systemInfo.Architecture, // 系统架构
			StartTime:           systemInfo.BootTime,     // 系统启动时间
		},
		NetworkList: networkList, // 节点的网卡列表信息
	}

	// 发送注册请求到服务器
	_, err = nc.client.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("注册请求失败: %v", err)
	}

	// 注册成功，记录日志
	logrus.Infof("节点注册成功，上报信息: %+v", req)
	return nil
}
