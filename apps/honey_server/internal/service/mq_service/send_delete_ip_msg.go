package mq_service

// File: service/mq_service/send_delete_ip_msg.go
// Description: 负责构建并发送删除IP相关的消息到RabbitMQ，处理批量IP删除任务的消息序列化及发布逻辑

import (
	"encoding/json"
	"honey_server/internal/global"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// DeleteIPRequest 删除IP的消息请求结构体
type DeleteIPRequest struct {
	IpList []IpInfo `json:"ipList"` // 待删除的IP信息列表
	LogID  string   `json:"logID"`  // 日志ID，用于追踪该批量删除任务的日志
}

// IpInfo 单个IP的详细信息结构体
type IpInfo struct {
	HoneyIPID uint   `json:"honeyIpID"` // 关联的诱捕IP记录ID
	IP        string `json:"ip"`        // 待删除的IP地址
	Network   string `json:"network"`   // 该IP所属的网络接口名称
	IsTan     bool   `json:"isTan"`     // 是否是探针IP
}

// SendDeleteIPMsg 发送批量删除IP的消息到指定节点的消息队列
func SendDeleteIPMsg(nodeUID string, req DeleteIPRequest) {
	// 将批量删除IP的请求参数序列化为JSON字节数据（消息体）
	byteData, _ := json.Marshal(req)
	cfg := global.Config.MQ // 获取全局MQ配置

	// 发布消息到RabbitMQ的删除IP业务交换器
	err := global.Queue.Publish(
		cfg.DeleteIpExchangeName, // 目标交换器名称（删除IP业务对应的交换器）
		nodeUID,                  // 路由键（目标节点UID，确保消息路由到指定节点队列）
		false,                    // mandatory：消息无法路由时是否返回给生产者（false不返回）
		false,                    // immediate：无消费者时是否立即返回（false不立即返回）
		amqp.Publishing{
			ContentType: "text/plain", // 消息内容类型
			Body:        byteData,     // 序列化后的消息体
		})

	// 消息发送结果日志记录
	if err != nil {
		logrus.Errorf("删除IP消息发送失败 %s %s", err, string(byteData))
	} else {
		logrus.Infof("删除IP消息发送成功 %s ", string(byteData))
	}
}
