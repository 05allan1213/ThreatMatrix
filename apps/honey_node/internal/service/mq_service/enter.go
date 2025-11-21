package mq_service

// File: service/mq_service/enter.go
// Description: 负责注册各类业务消息的消费处理器，完成交换器声明、队列绑定及消息消费逻辑处理

import (
	"fmt"
	"honey_node/internal/global"

	"github.com/sirupsen/logrus"
)

// Run 启动消息队列消费服务
func Run() {
	cfg := global.Config.MQ
	// 启动创建IP消息的消费协程
	go register(cfg.CreateIpExchangeName, CreateIpExChange)
	// 启动删除IP消息的消费协程
	go register(cfg.DeleteIpExchangeName, DeleteIpExChange)
	// 启动绑定端口消息的消费协程
	go register(cfg.BindPortExchangeName, BindPortExChange)
}

// register 注册单个业务的消息消费处理器
func register(exChangeName string, fun func(msg string) error) {
	// 声明与生产者一致的交换器（确保消费端与生产端交换器配置匹配）
	err := global.Queue.ExchangeDeclare(
		exChangeName, // 交换器名称（需与生产者一致）
		"direct",     // 交换器类型：direct（直接交换器，按路由键精准匹配）
		true,         // 持久化：交换器重启后不丢失
		false,        // 自动删除：无队列绑定时不自动删除
		false,        // 内部交换器：否（允许接收客户端消息）
		false,        // 非阻塞：否
		nil,          // 额外参数：无
	)
	if err != nil {
		logrus.Fatalf("%s 声明交换器失败 %s", exChangeName, err)
	}

	cf := global.Config
	// 声明队列（每个节点唯一队列，避免消息重复消费）
	queue, err := global.Queue.QueueDeclare(
		fmt.Sprintf("exChangeName_%s_queue", cf.System.Uid), // 队列名称：结合系统UID保证唯一性
		true,  // 持久化队列：队列重启后不丢失
		false, // 自动删除：否
		false, // 排他性：否（允许其他连接访问）
		false, // 非阻塞：否
		nil,   // 额外参数：无
	)
	if err != nil {
		logrus.Fatalf("声明队列失败 %s", err)
	}

	// 将队列绑定到交换器（通过路由键匹配消息）
	err = global.Queue.QueueBind(
		queue.Name,    // 队列名称
		cf.System.Uid, // 绑定键（路由键）：与生产者发送消息的路由键一致，确保消息精准路由到当前节点队列
		exChangeName,  // 交换器名称
		false,         // 非阻塞：否
		nil,           // 额外参数：无
	)
	if err != nil {
		logrus.Fatalf("绑定队列失败 %s", err)
	}

	// 注册消费者，开始监听队列消息（关闭自动确认，手动控制消息ack）
	msgs, err := global.Queue.Consume(
		queue.Name, // 消费的队列名称
		"",         // 消费者标签：空（使用默认标识）
		false,      // 自动确认：否（需手动ack/nack）
		false,      // 排他性：否
		false,      // 非本地消费者：否（接收所有发送到队列的消息）
		false,      // 非阻塞：否
		nil,        // 额外参数：无
	)
	if err != nil {
		logrus.Fatalf("注册消费者失败 %s", err)
	}

	// 循环消费队列中的消息
	for d := range msgs {
		// 调用业务处理函数处理消息内容
		err = fun(string(d.Body))
		if err == nil {
			// 消息处理成功：手动确认消息（false表示仅确认当前消息）
			d.Ack(false)
			continue
		}
		// 消息处理失败：拒绝消息并重新入队（false仅拒绝当前消息，true重新入队）
		d.Nack(false, true)
	}
}
