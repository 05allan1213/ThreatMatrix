package mq_service

// File: service/mq_service/bind_port_exchange.go
// Description: 处理端口绑定的消息，解析端口转发配置并启动TCP隧道转发服务，实现本地端口到目标地址的映射

import (
	"encoding/json"
	"fmt"
	"honey_node/internal/global"
	"honey_node/internal/models"
	"honey_node/internal/service/port_service"

	"github.com/sirupsen/logrus"
)

// BindPortRequest 端口绑定的消息请求结构体
type BindPortRequest struct {
	IP       string     `json:"ip"`       // 绑定端口的目标IP地址（本地监听的IP）
	PortList []PortInfo `json:"portList"` // 端口转发配置列表（每个项对应一个端口的转发规则）
	LogID    string     `json:"logID"`    // 日志ID，用于关联整个端口绑定流程的日志
}

// PortInfo 单个端口的转发配置结构体
type PortInfo struct {
	IP       string `json:"ip"`       // 本地监听的源IP地址
	Port     int    `json:"port"`     // 本地监听的源端口号
	DestIP   string `json:"destIP"`   // 转发的目标IP地址
	DestPort int    `json:"destPort"` // 转发的目标端口号
}

// LocalAddr 拼接本地监听的地址字符串（IP:Port）
// 用于port_service建立本地监听时的地址参数
func (p PortInfo) LocalAddr() string {
	return fmt.Sprintf("%s:%d", p.IP, p.Port)
}

// TargetAddr 拼接转发的目标地址字符串（DestIP:DestPort）
// 用于port_service转发数据时的目标地址参数
func (p PortInfo) TargetAddr() string {
	return fmt.Sprintf("%s:%d", p.DestIP, p.DestPort)
}

// BindPortExChange 处理端口绑定的消息消费逻辑
// 解析端口转发配置，为每个端口启动独立的TCP隧道转发服务，实现本地端口到目标地址的映射
func BindPortExChange(msg string) error {
	// 记录接收到的端口绑定消息，便于调试
	logrus.Infof("端口绑定消息 %#v", msg)

	var req BindPortRequest
	// 将JSON消息体解析为PortBindRequest结构体
	if err := json.Unmarshal([]byte(msg), &req); err != nil {
		logrus.Errorf("JSON解析失败: %v, 消息: %s", err, msg)
		return nil // 解析失败返回nil，避免消息重复处理
	}

	// 把这个IP上的服务全部停掉
	port_service.CloseIpTunnel(req.IP)

	// 遍历端口转发配置列表，为每个端口启动独立的转发服务
	for _, port := range req.PortList {
		global.DB.Create(&models.PortModel{
			TargetAddr: port.TargetAddr(),
			LocalAddr:  port.LocalAddr(),
		})
		// 使用goroutine异步启动每个端口的转发，避免阻塞消息处理
		go func(port PortInfo) {
			// 调用port_service的Tunnel方法，建立本地端口到目标地址的TCP转发隧道
			err := port_service.Tunnel(port.LocalAddr(), port.TargetAddr())
			if err != nil {
				logrus.Errorf("端口绑定失败 %s", err)
			}
		}(port)
	}

	return nil
}
