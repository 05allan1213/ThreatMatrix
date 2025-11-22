package mq_service

// File: service/mq_service/create_ip_exchange.go
// Description: 创建IP消息处理器，集成ARP检测避免IP冲突、资源自动清理、模块化命令执行及统一状态上报机制

import (
	"context"
	"encoding/json"
	"fmt"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/cmd"
	"net"
	"strings"

	"github.com/j-keck/arping"
	"github.com/sirupsen/logrus"
)

// CreateIPRequest 创建IP的消息请求结构体
type CreateIPRequest struct {
	HoneyIPID uint   `json:"honeyIpID"` // 关联的诱捕IP记录ID，用于生成唯一的macvlan接口名称
	IP        string `json:"ip"`        // 待配置的IP地址（需先通过ARP检测确认未被占用）
	Mask      int8   `json:"mask"`      // 子网掩码位数（如24表示255.255.255.0）
	Network   string `json:"network"`   // 基于哪个物理/主网络接口创建macvlan子接口（如ens33）
	LogID     string `json:"logID"`     // 日志追踪ID，用于关联整个创建流程的日志
	IsTan     bool   `json:"isTan"`     // 是否是探针ip
}

// CreateIpExChange 创建IP消息处理函数
func CreateIpExChange(msg string) error {
	var req CreateIPRequest
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		logrus.Errorf("JSON解析失败: %v, 消息: %s", err, msg)
		return nil // 解析失败返回nil，避免消息重复处理
	}

	// 探针ip直接上报状态
	if req.IsTan {
		mac, _ := getMACAddress(req.Network)
		return reportStatus(req.HoneyIPID, req.Network, mac, "")
	}

	// 记录创建请求的详细上下文日志，便于问题排查
	global.Log.WithFields(logrus.Fields{
		"honeyIPID": req.HoneyIPID,
		"ip":        req.IP,
		"mask":      req.Mask,
		"network":   req.Network,
		"logID":     req.LogID,
	}).Info("开始处理创建IP请求")

	// ARP预检测：检查目标IP是否已被局域网内其他设备占用
	_mac, _, err := arping.PingOverIfaceByName(net.ParseIP(req.IP), req.Network)
	if err == nil {
		err = fmt.Errorf("创建诱捕ip失败 ip已存在 ip %s mac %s", req.IP, _mac.String())
		logrus.Error(err)
		return reportStatus(req.HoneyIPID, "", _mac.String(), err.Error()) // 直接上报冲突状态
	}

	linkName := fmt.Sprintf("hy_%d", req.HoneyIPID) // 生成唯一的macvlan接口名称

	// 定义资源清理函数：操作失败时自动清理已创建的网络接口，避免资源泄漏
	cleanup := func() {
		if err := cmd.Cmd(fmt.Sprintf("ip link delete %s", linkName)); err != nil {
			logrus.Errorf("清理失败，删除网络接口 %s 时出错: %v", linkName, err)
		}
	}

	// 分步执行网络配置命令，每步失败均触发资源清理并上报状态
	if err := createMacVlanInterface(linkName, req.Network); err != nil {
		logrus.Errorf("创建macvlan接口失败: %v", err)
		cleanup()
		return reportStatus(req.HoneyIPID, linkName, "", err.Error())
	}

	if err := setInterfaceUp(linkName); err != nil {
		logrus.Errorf("启用网络接口失败: %v", err)
		cleanup()
		return reportStatus(req.HoneyIPID, linkName, "", err.Error())
	}

	if err := addIPAddress(linkName, req.IP, req.Mask); err != nil {
		logrus.Errorf("添加IP地址失败: %v", err)
		cleanup()
		return reportStatus(req.HoneyIPID, linkName, "", err.Error())
	}

	mac, err := getMACAddress(linkName)
	if err != nil {
		logrus.Errorf("获取MAC地址失败: %v", err)
		cleanup()
		return reportStatus(req.HoneyIPID, linkName, "", err.Error())
	}

	// 所有步骤执行成功，上报最终状态
	return reportStatus(req.HoneyIPID, linkName, mac, "")
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

// addIPAddress 为网络接口配置IP地址和子网掩码
func addIPAddress(linkName, ip string, mask int8) error {
	cmdStr := fmt.Sprintf("ip addr add %s/%d dev %s", ip, mask, linkName)
	if err := cmd.Cmd(cmdStr); err != nil {
		return fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return nil
}

// getMACAddress 获取网络接口的MAC地址
func getMACAddress(linkName string) (string, error) {
	cmdStr := fmt.Sprintf("ip link show %s | awk '/link\\/ether/ {print $2}'", linkName)
	mac, err := cmd.Command(cmdStr)
	if err != nil {
		return "", fmt.Errorf("执行命令失败 [%s]: %w", cmdStr, err)
	}
	return strings.TrimSpace(mac), nil // 去除MAC地址前后的空白字符
}

// reportStatus 统一状态上报函数
func reportStatus(honeyIPID uint, network, mac, errMsg string) error {
	response, err := global.GrpcClient.StatusCreateIP(context.Background(), &node_rpc.StatusCreateIPRequest{
		HoneyIPID: uint32(honeyIPID),
		ErrMsg:    errMsg,
		Network:   network,
		Mac:       mac,
	})

	if err != nil {
		logrus.Errorf("上报管理状态失败: %v", err)
		return err
	}

	global.Log.WithFields(logrus.Fields{
		"honeyIPID": honeyIPID,
		"network":   network,
		"mac":       mac,
		"errMsg":    errMsg,
	}).Infof("上报管理状态成功: %v", response)

	return nil
}
