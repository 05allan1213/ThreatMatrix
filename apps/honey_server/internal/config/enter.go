package config

// File: config/enter.go
// Description: 定义诱捕服务所需的配置和相关方法。

import "fmt"

// 统一管理配置结构体
type Config struct {
	DB        DB       `yaml:"db"`
	Logger    Logger   `yaml:"logger"`
	Redis     Redis    `yaml:"redis"`
	System    System   `yaml:"system"`
	Jwt       Jwt      `yaml:"jwt"`
	WhiteList []string `yaml:"whiteList"`
	MQ        MQ       `yaml:"mq"`
}

// 数据库配置
type DB struct {
	DbName          string `yaml:"db_name"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime int    `yaml:"connMaxLifetime"`
}

// DSN 构建数据库连接字符串
func (cfg DB) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)
}

// 日志配置
type Logger struct {
	Format  string `yaml:"format"`
	Level   string `yaml:"level"`
	AppName string `yaml:"appName"`
}

// Redis 配置
type Redis struct {
	Addr     string
	Password string
	DB       int
}

// 系统配置
type System struct {
	WebAddr  string `yaml:"webAddr"`
	GrpcAddr string `yaml:"grpcAddr"`
	Mode     string `yaml:"mode"`
}

// Jwt 配置
type Jwt struct {
	Expires int    `yaml:"expires"` // 过期时间，单位: 秒
	Issuer  string `yaml:"issuer"`  // 签发者
	Secret  string `yaml:"secret"`  // 密钥
}

// rabbitMQ 配置
type MQ struct {
	User                 string `yaml:"user"`                 // RabbitMQ 用户名
	Password             string `yaml:"password"`             // RabbitMQ 密码
	Host                 string `yaml:"host"`                 // RabbitMQ 主机地址
	Port                 int    `yaml:"port"`                 // RabbitMQ 端口号
	CreateIpExchangeName string `yaml:"createIpExchangeName"` // 创建IP交换机名称
	DeleteIpExchangeName string `yaml:"deleteIpExchangeName"` // 删除IP交换机名称
	BindPortExchangeName string `yaml:"bindPortExchangeName"` // 绑定端口交换机名称
	Ssl                  bool   `yaml:"ssl"`                  // 是否启用SSL
	ClientCertificate    string `yaml:"clientCertificate"`    // 客户端证书路径
	ClientKey            string `yaml:"clientKey"`            // 客户端密钥路径
	CaCertificate        string `yaml:"caCertificate"`        // CA证书路径
}

// 构造并返回RabbitMQ服务器的连接地址
func (m MQ) Addr() string {
	return fmt.Sprintf("amqps://%s:%s@%s:%d/",
		m.User,
		m.Password,
		m.Host,
		m.Port,
	)
}
