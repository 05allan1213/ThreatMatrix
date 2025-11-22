package mq_service

// File: service/mq_service/delete_ip_exchange.go
// Description: 处理删除诱捕IP的消息，执行网络接口删除命令并批量上报删除状态到服务端

import (
	"context"
	"encoding/json"
	"fmt"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/cmd"

	"github.com/sirupsen/logrus"
)

// DeleteIPRequest 删除诱捕IP的消息请求结构体
type DeleteIPRequest struct {
	IpList []IpInfo `json:"ipList"` // 待删除的诱捕IP信息列表
	LogID  string   `json:"logID"`  // 日志ID，用于关联整个批量删除流程的日志
}

// IpInfo 单个诱捕IP的详细信息结构体
type IpInfo struct {
	HoneyIPID uint   `json:"honeyIpID"` // 关联的诱捕IP记录ID
	IP        string `json:"ip"`        // 待删除的IP地址
	Network   string `json:"network"`   // 对应的macvlan网络接口名称（如hy_12）
}

// DeleteIpExChange 处理删除诱捕IP的消息消费逻辑
func DeleteIpExChange(msg string) error {
	var req DeleteIPRequest
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		logrus.Errorf("JSON解析失败: %v, 消息: %s", err, msg)
		return nil // 解析失败返回nil，避免消息重复处理
	}

	// 记录批量删除请求的上下文日志，便于追踪处理流程
	global.Log.WithFields(logrus.Fields{
		"req": req,
	}).Infof("删除诱捕ip")

	var idList []uint32 // 收集待上报的诱捕IPID列表
	// 遍历IP列表，执行网络接口删除命令
	for _, info := range req.IpList {
		// 删除对应的macvlan网络接口（info.Network为接口名称）
		cmd.Cmd(fmt.Sprintf("ip link del %s", info.Network))
		idList = append(idList, uint32(info.HoneyIPID)) // 收集ID用于状态上报
	}

	// 上报批量删除状态到服务端
	reportDeleteIPStatus(idList)
	return nil
}

// reportDeleteIPStatus 批量上报删除诱捕IP的状态到服务端
func reportDeleteIPStatus(honeyIPIDList []uint32) error {
	response, err := global.GrpcClient.StatusDeleteIP(context.Background(), &node_rpc.StatusDeleteIPRequest{
		HoneyIPIDList: honeyIPIDList, // 批量删除的IPID列表
	})

	if err != nil {
		logrus.Errorf("上报管理状态失败: %v", err)
		return err
	}

	// 记录状态上报成功的日志，包含ID列表便于核对
	global.Log.WithFields(logrus.Fields{
		"honeyIPIDList": honeyIPIDList,
	}).Infof("上报管理状态成功: %v", response)

	return nil
}
