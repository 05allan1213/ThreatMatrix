package config

import "fmt"

// File: config/enter.go
// Description: 定义诱捕服务所需的配置和相关方法。

// 统一管理配置结构体
type Config struct {
	Logger            Logger   `yaml:"logger"`
	System            System   `yaml:"system"`
	FilterNetworkList []string `yaml:"filterNetworkList"`
	MQ                MQ       `yaml:"mq"`
	DB                DB       `yaml:"db"`
}

// 日志配置
type Logger struct {
	Format  string `yaml:"format"`
	Level   string `yaml:"level"`
	AppName string `yaml:"appName"`
}

// 系统配置
type System struct {
	GrpcManageAddr string `yaml:"grpcManageAddr"`
	Network        string `yaml:"network"`
	Uid            string `yaml:"uid"`
}

// rabbitMQ 配置
type MQ struct {
	User                 string `yaml:"user"`
	Password             string `yaml:"password"`
	Host                 string `yaml:"host"`
	Port                 int    `yaml:"port"`
	CreateIpExchangeName string `yaml:"createIpExchangeName"`
	DeleteIpExchangeName string `yaml:"deleteIpExchangeName"`
	BindPortExchangeName string `yaml:"bindPortExchangeName"`
	Ssl                  bool   `yaml:"ssl"`
	ClientCertificate    string `yaml:"clientCertificate"`
	ClientKey            string `yaml:"clientKey"`
	CaCertificate        string `yaml:"caCertificate"`
}

// 获取rabbitMQ连接地址
func (m MQ) Addr() string {
	// 如果启用ssl，则返回ssl连接地址
	if m.Ssl {
		return fmt.Sprintf("amqps://%s:%s@%s:%d/",
			m.User,
			m.Password,
			m.Host,
			m.Port,
		)
	}
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		m.User,
		m.Password,
		m.Host,
		m.Port,
	)
}

// 数据库配置
type DB struct {
	DbName          string `yaml:"db_name"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}
