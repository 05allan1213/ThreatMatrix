package mq_service

// File: service/mq_service/create_ip_exchange.go
// Description: 负责处理创建IP的消息

import (
	"encoding/json"
	"fmt"
	"honey_node/internal/utils/cmd"
	"strings"

	"github.com/sirupsen/logrus"
)

// CreateIPRequest 创建IP的消息请求结构体
type CreateIPRequest struct {
	HoneyIPID uint   `json:"honeyIpID"` // 关联的诱捕IP记录ID，用于生成唯一的接口名称
	IP        string `json:"ip"`        // 要配置的IP地址
	Mask      int8   `json:"mask"`      // 子网掩码位数（如24表示255.255.255.0）
	Network   string `json:"network"`   // 基于哪个物理/主网络接口创建macvlan子接口
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

	// 生成macvlan子接口名称（格式：hy_+诱捕IPID，确保唯一性）
	linkName := fmt.Sprintf("hy_%d", req.HoneyIPID)

	// 1. 创建macvlan子接口：基于指定的主接口（req.Network）创建桥接模式的macvlan接口
	cmd.Cmd(fmt.Sprintf("ip link add %s link %s type macvlan mode bridge", linkName, req.Network))
	// 2. 启用该macvlan子接口
	cmd.Cmd(fmt.Sprintf("ip link set %s up", linkName))
	// 3. 为macvlan子接口配置IP地址和子网掩码
	cmd.Cmd(fmt.Sprintf("ip addr add %s/%d dev %s", req.IP, req.Mask, linkName))

	// 获取macvlan子接口的MAC地址：通过ip命令结合awk过滤提取
	mac, err := cmd.Command(fmt.Sprintf("ip link show %s | awk '/link\\/ether/ {print $2}'", linkName))
	if err != nil {
		logrus.Errorf("获取macvlan子接口MAC地址失败 %s", err)
		return nil // 获取MAC地址失败时返回nil，避免消息重复处理
	}
	fmt.Println("mac: ", strings.TrimSpace(mac)) // 打印并去除MAC地址前后的空白字符

	/* 命令示例参考：
	ip link add mc_12 link ens33 type macvlan mode bridge  // 创建macvlan接口
	ip link set mc_12 up                                   // 启用接口
	ip addr add 192.168.80.166/24 dev mc_12               // 配置IP
	*/

	// 调grpc方法，上报状态（此处为预留逻辑，实际需补充grpc调用代码）

	return nil
}
