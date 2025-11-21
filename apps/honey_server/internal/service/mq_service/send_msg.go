package mq_service

// File: service/mq_service/send_msg.go
// Description: 负责构建并发送创建IP相关的消息到RabbitMQ，处理消息序列化及发布逻辑

import (
	"encoding/json"
	"honey_server/internal/global"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// CreateIPRequest 创建IP的消息请求结构体
type CreateIPRequest struct {
	HoneyIPID uint   `json:"honeyIpID"` // 关联的诱捕IP记录ID
	IP        string `json:"ip"`        // 要创建的IP地址
	Mask      int8   `json:"mask"`      // 子网掩码位数
	Network   string `json:"network"`   // 基于哪个网络接口创建
	LogID     string `json:"logID"`     // 日志ID，用于追踪该任务的日志
}

// SendCreateIPMsg 发送创建IP的消息到指定节点的消息队列
func SendCreateIPMsg(nodeUID string, req CreateIPRequest) {
	// 将请求参数序列化为JSON字节数据（消息体）
	byteData, _ := json.Marshal(req)
	cfg := global.Config.MQ // 获取全局MQ配置

	// 发布消息到RabbitMQ交换器
	err := global.Queue.Publish(
		cfg.CreateIpExchangeName, // 目标交换器名称（创建IP业务对应的交换器）
		nodeUID,                  // 路由键（目标节点UID，确保消息路由到指定节点队列）
		false,                    // mandatory：消息无法路由时是否返回给生产者（false不返回）
		false,                    // immediate：无消费者时是否立即返回（false不立即返回）
		amqp.Publishing{
			ContentType: "text/plain", // 消息内容类型
			Body:        byteData,     // 序列化后的消息体
		})

	// 消息发送结果日志记录
	if err != nil {
		logrus.Errorf("消息发送失败 %s %s", err, string(byteData))
	} else {
		logrus.Infof("消息发送成功 %s ", string(byteData))
	}
}
