package core

// File: core/mq.go
// Description: 负责初始化RabbitMQ消息队列连接

import (
	"crypto/tls"
	"crypto/x509"
	"honey_node/internal/global"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// InitMQ 初始化RabbitMQ连接并返回消息通道
func InitMQ() *amqp.Channel {
	cfg := global.Config.MQ // 获取全局配置中的MQ配置项
	var conn *amqp.Connection
	var err error

	// 判断是否启用SSL/TLS加密连接
	if cfg.Ssl {
		// 1. 加载客户端证书和私钥（用于双向认证，服务端验证客户端身份）
		cert, err := tls.LoadX509KeyPair(cfg.ClientCertificate, cfg.ClientKey)
		if err != nil {
			logrus.Fatalf("加载客户端证书失败: %v", err)
		}

		// 2. 加载CA根证书（用于验证服务端证书的合法性）
		caCert, err := os.ReadFile(cfg.CaCertificate)
		if err != nil {
			logrus.Fatalf("读取CA证书失败: %v", err)
		}
		caCertPool := x509.NewCertPool()      // 创建CA证书池
		caCertPool.AppendCertsFromPEM(caCert) // 将CA证书添加到信任池

		// 3. 配置TLS连接参数
		tlsConfig := &tls.Config{
			Certificates:       []tls.Certificate{cert}, // 客户端证书链（双向认证必需）
			RootCAs:            caCertPool,              // 信任的CA根证书池
			InsecureSkipVerify: false,                   // 禁止跳过服务端证书验证，确保安全性
		}

		// 通过TLS加密方式连接RabbitMQ
		conn, err = amqp.DialTLS(cfg.Addr(), tlsConfig)
		if err != nil {
			logrus.Fatalf("无法连接到 RabbitMQ: %v", err)
		}
	} else {
		// 使用非加密方式连接RabbitMQ
		conn, err = amqp.Dial(cfg.Addr())
	}

	// 连接失败则终止程序并记录日志
	if err != nil {
		logrus.Fatalf("无法连接到 RabbitMQ: %v", err)
	}

	// 创建RabbitMQ消息通道（Channel）
	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("无法打开通道: %v", err)
	}

	return ch // 返回创建好的消息通道
}
