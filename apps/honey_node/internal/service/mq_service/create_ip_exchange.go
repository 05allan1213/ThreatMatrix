package mq_service

// File: service/mq_service/create_ip_exchange.go
// Description: 处理创建IP的消息，执行macvlan接口创建、IP配置命令，获取MAC地址并通过gRPC上报创建状态

import (
	"context"
	"encoding/json"
	"fmt"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/cmd"
	"strings"

	"github.com/sirupsen/logrus"
)

// CreateIPRequest 创建IP的消息请求结构体
type CreateIPRequest struct {
	HoneyIPID uint   `json:"honeyIpID"` // 关联的诱捕IP记录ID，用于生成唯一的macvlan接口名称
	IP        string `json:"ip"`        // 要配置的IP地址
	Mask      int8   `json:"mask"`      // 子网掩码位数（如24表示255.255.255.0）
	Network   string `json:"network"`   // 基于哪个物理/主网络接口创建macvlan子接口（如ens33）
	LogID     string `json:"logID"`     // 日志ID，用于追踪该创建任务的日志
}

// CreateIpExChange 处理创建IP的消息消费逻辑
func CreateIpExChange(msg string) error {
	var req CreateIPRequest
	// 将JSON消息体解析为CreateIPRequest结构体
	err := json.Unmarshal([]byte(msg), &req)
	if err != nil {
		logrus.Errorf("json解析失败 %s %s", err, msg)
		return nil // 解析失败时返回nil，避免消息重复处理
	}

	var errorMsg string
	// 生成唯一的macvlan子接口名称（格式：hy_+诱捕IPID）
	linkName := fmt.Sprintf("hy_%d", req.HoneyIPID)

	// 1. 创建macvlan子接口：基于指定主接口创建桥接模式的macvlan接口
	err = cmd.Cmd(fmt.Sprintf("ip link add %s link %s type macvlan mode bridge", linkName, req.Network))
	if err != nil {
		errorMsg = err.Error() // 捕获命令执行错误
	}

	// 2. 启用macvlan子接口
	err = cmd.Cmd(fmt.Sprintf("ip link set %s up", linkName))
	if err != nil {
		errorMsg = err.Error() // 捕获命令执行错误
	}

	// 3. 为macvlan子接口配置IP地址和子网掩码
	err = cmd.Cmd(fmt.Sprintf("ip addr add %s/%d dev %s", req.IP, req.Mask, linkName))
	if err != nil {
		errorMsg = err.Error() // 捕获命令执行错误
	}

	// 获取macvlan子接口的MAC地址：通过ip命令结合awk过滤提取
	mac, err := cmd.Command(fmt.Sprintf("ip link show %s | awk '/link\\/ether/ {print $2}'", linkName))
	if err != nil {
		errorMsg = err.Error() // 捕获MAC地址获取错误
	}
	mac = strings.TrimSpace(mac) // 去除MAC地址前后的空白字符

	/* 命令示例参考：
	ip link add mc_12 link ens33 type macvlan mode bridge  // 创建macvlan接口
	ip link set mc_12 up                                   // 启用接口
	ip addr add 192.168.80.166/24 dev mc_12               // 配置IP
	*/

	// 通过gRPC调用上报创建状态到服务端
	response, err := global.GrpcClient.StatusCreateIP(context.Background(), &node_rpc.StatusCreateIPRequest{
		HoneyIPID: uint32(req.HoneyIPID), // 诱捕IPID
		ErrMsg:    errorMsg,              // 执行过程中的错误信息（空表示成功）
		Network:   linkName,              // 创建的macvlan接口名称
		Mac:       mac,                   // 获取到的接口MAC地址
	})
	if err != nil {
		logrus.Errorf("上报管理状态失败 %s", err)
		return err
	}
	logrus.Infof("上报管理状态成功 %v", response)
	return nil
}
