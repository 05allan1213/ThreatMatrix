package mq_service

// File: service/mq_service/create_ip_exchange.go
// Description: 创建IP交换机

import "fmt"

// CreateIpExChange 创建IP交换机
func CreateIpExChange(msg string) error {
	fmt.Println("消息", msg)
	return nil
}
