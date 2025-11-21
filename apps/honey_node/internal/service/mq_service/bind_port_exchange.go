package mq_service

import "fmt"

// File:service/mq_service/bind_port_exchange.go
// Description: 绑定端口交换机

// BindPortExChange 绑定端口交换机
func BindPortExChange(msg string) error {
	fmt.Println("端口绑定消息", msg)
	return nil
}
