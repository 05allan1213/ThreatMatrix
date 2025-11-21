package mq_service

// File: service/mq_service/send_bind_port_msg.go
// Description: 负责构建并发送端口绑定相关的消息到RabbitMQ，处理端口转发配置任务的消息序列化及发布逻辑

import (
	"encoding/json"
	"honey_server/internal/global"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// BindPortRequest 端口绑定的消息请求结构体
type BindPortRequest struct {
	IP       string     `json:"ip"`       // 绑定端口的目标IP地址
	PortList []PortInfo `json:"portList"` // 端口绑定配置列表
	LogID    string     `json:"logID"`    // 日志ID，用于追踪该端口绑定任务的日志
}

// PortInfo 单个端口的绑定信息结构体
type PortInfo struct {
	IP       string `json:"ip"`       // 源IP地址
	Port     int    `json:"port"`     // 源端口号
	DestIP   string `json:"destIP"`   // 转发目标IP地址
	DestPort int    `json:"destPort"` // 转发目标端口号
}

// SendBindPortMsg 发送端口绑定的消息到指定节点的消息队列
func SendBindPortMsg(nodeUID string, req BindPortRequest) {
	// 将端口绑定的请求参数序列化为JSON字节数据（消息体）
	byteData, _ := json.Marshal(req)
	cfg := global.Config.MQ // 获取全局MQ配置

	// 发布消息到RabbitMQ的端口绑定业务交换器
	err := global.Queue.Publish(
		cfg.BindPortExchangeName, // 目标交换器名称（端口绑定业务对应的交换器）
		nodeUID,                  // 路由键（目标节点UID，确保消息路由到指定节点队列）
		false,                    // mandatory：消息无法路由时是否返回给生产者（false不返回）
		false,                    // immediate：无消费者时是否立即返回（false不立即返回）
		amqp.Publishing{
			ContentType: "text/plain", // 消息内容类型
			Body:        byteData,     // 序列化后的消息体
		})

	// 消息发送结果日志记录
	if err != nil {
		logrus.Errorf("端口绑定消息发送失败 %s %s", err, string(byteData))
	} else {
		logrus.Infof("端口绑定消息发送成功 %s ", string(byteData))
	}
}
