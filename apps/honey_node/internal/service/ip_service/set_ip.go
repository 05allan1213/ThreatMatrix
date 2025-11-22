package ip_service

// File: service/ip_service/set_ip.go
// Description: 提供macvlan接口的创建、MAC地址配置、IP地址分配等功能，集成资源自动清理机制确保操作原子性

import (
	"fmt"
	"honey_node/internal/utils/cmd"
	"strings"

	"github.com/sirupsen/logrus"
)

// SetIpRequest 创建macvlan网络接口的请求结构体
type SetIpRequest struct {
	Ip       string `json:"ip"`       // 待配置的IP地址
	Mask     int8   `json:"mask"`     // 子网掩码位数（如24表示255.255.255.0）
	LinkName string `json:"linkName"` // 要创建的macvlan子接口名称（如hy_12）
	Network  string `json:"network"`  // 基于哪个物理/主网络接口创建（如ens33）
	Mac      string `json:"mac"`      // 接口的MAC地址（为空时由系统自动分配，非空时手动指定）
}

// SetIp 创建并配置macvlan网络接口
func SetIp(req SetIpRequest) (mac string, err error) {
	linkName := req.LinkName

	// 定义资源清理函数：操作失败时删除已创建的网络接口，避免资源泄漏
	cleanup := func() {
		if err := cmd.Cmd(fmt.Sprintf("ip link delete %s", linkName)); err != nil {
			logrus.Errorf("清理失败，删除网络接口 %s 时出错: %v", linkName, err)
		}
	}

	// 1. 创建macvlan子接口（基于指定的主网卡）
	if err = createMacVlanInterface(linkName, req.Network); err != nil {
		logrus.Errorf("创建macvlan接口失败: %v", err)
		cleanup() // 创建失败时清理资源
		return
	}

	// 2. 若指定了MAC地址，则手动设置（用于固定接口标识）
	if req.Mac != "" {
		err = setInterfaceMac(linkName, req.Mac)
		if err != nil {
			logrus.Errorf("设置mac失败: %v", err)
			cleanup() // MAC设置失败时清理资源
			return
		}
	}

	// 3. 启用macvlan子接口（使其处于工作状态）
	if err = setInterfaceUp(linkName); err != nil {
		logrus.Errorf("启用网络接口失败: %v", err)
		cleanup() // 启用失败时清理资源
		return
	}

	// 4. 为接口配置IP地址和子网掩码
	if err = addIPAddress(linkName, req.Ip, req.Mask); err != nil {
		logrus.Errorf("添加IP地址失败: %v", err)
		cleanup() // IP配置失败时清理资源
		return
	}

	// 5. 若未指定MAC地址，则从系统中获取自动分配的MAC地址
	if req.Mac == "" {
		req.Mac, err = GetMACAddress(linkName)
		if err != nil {
			cleanup() // 获取MAC失败时清理资源
			return
		}
	}

	return req.Mac, nil // 返回最终的MAC地址（手动指定或自动分配）
}

// createMacVlanInterface 创建macvlan子接口
func createMacVlanInterface(linkName, network string) error {
	cmdStr := fmt.Sprintf("ip link add %s link %s type macvlan mode bridge", linkName, network)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err) // 包装错误信息，保留原始错误链
	}
	return nil
}

// setInterfaceUp 启用网络接口
func setInterfaceUp(linkName string) error {
	cmdStr := fmt.Sprintf("ip link set %s up", linkName)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return nil
}

// setInterfaceMac 手动设置网络接口的MAC地址
func setInterfaceMac(linkName string, mac string) error {
	cmdStr := fmt.Sprintf("ip link set %s address %s", linkName, mac)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return nil
}

// addIPAddress 为网络接口配置IP地址和子网掩码
func addIPAddress(linkName, ip string, mask int8) error {
	cmdStr := fmt.Sprintf("ip addr add %s/%d dev %s", ip, mask, linkName)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return nil
}

// GetMACAddress 获取网络接口的MAC地址
func GetMACAddress(linkName string) (string, error) {
	cmdStr := fmt.Sprintf("ip link show %s | awk '/link\\/ether/ {print $2}'", linkName)
	mac, err := cmd.Command(cmdStr)
	if err != nil {
		return "", fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return strings.TrimSpace(mac), nil // 去除MAC地址前后的空白字符，确保格式整洁
}
