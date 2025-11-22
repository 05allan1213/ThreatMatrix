package mq_service

// File: service/mq_service/create_ip_exchange.go
// Description: 创建IP消息处理器，集成ARP检测避免IP冲突、资源自动清理、模块化命令执行及统一状态上报机制

import (
	"context"
	"encoding/json"
	"fmt"
	"honey_node/internal/global"
	"honey_node/internal/models"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/service/ip_service"
	"net"

	"github.com/j-keck/arping"
	"github.com/sirupsen/logrus"
)

// CreateIPRequest 创建IP的消息请求结构体
type CreateIPRequest struct {
	HoneyIPID uint   `json:"honeyIpID"` // 关联的诱捕IP记录ID
	IP        string `json:"ip"`        // 待配置的IP地址
	Mask      int8   `json:"mask"`      // 子网掩码位数
	Network   string `json:"network"`   // 基于哪个主网络接口创建
	LogID     string `json:"logID"`     // 日志追踪ID
	IsTan     bool   `json:"isTan"`     // 是否为探针IP（探针IP无需创建新接口，复用主接口）
}

// CreateIpExChange 处理创建IP的消息消费逻辑
func CreateIpExChange(msg string) error {
	var req CreateIPRequest
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		logrus.Errorf("JSON解析失败: %v, 消息: %s", err, msg)
		return nil // 解析失败返回nil，避免消息重复处理
	}

	// 探针IP处理逻辑：无需创建新接口，直接获取主接口MAC地址并上报
	if req.IsTan {
		mac, _ := ip_service.GetMACAddress(req.Network) // 获取主接口的MAC地址
		return reportStatus(req.HoneyIPID, req.Network, mac, "")
	}

	// 普通诱捕IP处理流程：记录上下文日志便于追踪
	global.Log.WithFields(logrus.Fields{
		"honeyIPID": req.HoneyIPID,
		"ip":        req.IP,
		"mask":      req.Mask,
		"network":   req.Network,
		"logID":     req.LogID,
	}).Info("开始处理创建IP请求")

	// ARP冲突检测：通过arping检查目标IP是否已被局域网内其他设备占用
	_mac, _, err := arping.PingOverIfaceByName(net.ParseIP(req.IP), req.Network)
	if err == nil {
		err = fmt.Errorf("创建诱捕ip失败 ip已存在 ip %s mac %s", req.IP, _mac.String())
		logrus.Error(err)
		return reportStatus(req.HoneyIPID, "", _mac.String(), err.Error()) // 上报IP冲突错误
	}

	// 生成唯一的macvlan子接口名称（格式：hy_+诱捕IPID）
	linkName := fmt.Sprintf("hy_%d", req.HoneyIPID)

	// 调用ip_service创建macvlan接口并配置IP
	mac, err := ip_service.SetIp(ip_service.SetIpRequest{
		Ip:       req.IP,
		Mask:     req.Mask,
		LinkName: linkName,
		Network:  req.Network,
	})
	if err != nil {
		return reportStatus(req.HoneyIPID, linkName, mac, err.Error()) // 创建失败时上报错误
	}

	// 配置持久化：将IP接口信息保存到数据库，支持程序重启后恢复配置
	global.DB.Create(&models.IpModel{
		Ip:       req.IP,
		Mask:     req.Mask,
		LinkName: linkName,
		Network:  req.Network,
		Mac:      mac,
	})

	// 所有步骤成功，上报创建成功状态
	return reportStatus(req.HoneyIPID, linkName, mac, "")
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

	logrus.WithFields(logrus.Fields{
		"honeyIPID": honeyIPID,
		"network":   network,
		"mac":       mac,
		"errMsg":    errMsg,
	}).Infof("上报管理状态成功: %v", response)

	return nil
}
