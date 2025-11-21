package mq_service

// File: service/mq_service/register_exchange.go
// Description: 负责RabbitMQ交换器（Exchange）的注册与声明

import (
	"honey_server/internal/global"

	"github.com/sirupsen/logrus"
)

// RegisterExChange 注册所有业务相关的RabbitMQ交换器
func RegisterExChange() {
	cfg := global.Config.MQ // 获取全局配置中的MQ配置项
	// 声明创建IP相关的交换器
	exchangeDeclare(cfg.CreateIpExchangeName)
	// 声明删除IP相关的交换器
	exchangeDeclare(cfg.DeleteIpExchangeName)
	// 声明绑定端口相关的交换器
	exchangeDeclare(cfg.BindPortExchangeName)
}

// exchangeDeclare 声明单个RabbitMQ交换器
func exchangeDeclare(name string) {
	// 调用RabbitMQ通道的ExchangeDeclare方法声明交换器
	err := global.Queue.ExchangeDeclare(
		name,     // 交换器名称
		"direct", // 交换器类型：direct（直接交换器），根据路由键完全匹配路由消息
		true,     // 持久化：交换器在服务器重启后仍然存在
		false,    // 自动删除：当所有绑定队列都解绑后，交换器不会自动删除
		false,    // 内部：是否为内部交换器（仅用于交换器间转发，不接收客户端消息）
		false,    // 非阻塞：是否阻塞等待声明完成
		nil,      // 额外参数：无特殊配置
	)
	// 声明失败则终止程序并记录日志
	if err != nil {
		logrus.Fatalf("声明交换器 %s 失败 %s", name, err)
		return
	}
	// 声明成功打印日志
	logrus.Infof("声明交换器 %s 成功", name)
}
