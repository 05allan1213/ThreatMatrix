package mq_service

import "fmt"

// File: service/mq_service/delete_ip_exchange.go
// Description: 删除IP交换机

// DeleteIpExChange 删除IP交换机
func DeleteIpExChange(msg string) error {
	fmt.Println("删除消息", msg)
	return nil
}
