package config

// File: config/enter.go
// Description: 定义诱捕服务所需的配置和相关方法。

import "fmt"

// 统一管理配置结构体
type Config struct {
	DB     DB     `yaml:"db"`
	Logger Logger `yaml:"logger"`
	Redis  Redis  `yaml:"redis"`
	System System `yaml:"system"`
	Jwt    Jwt    `yaml:"jwt"`
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
	WebAddr string `yaml:"webAddr"`
	Mode    string `yaml:"mode"`
}

// Jwt 配置
type Jwt struct {
	Expires int    `yaml:"expires"` // 过期时间，单位: 秒
	Issuer  string `yaml:"issuer"`  // 签发者
	Secret  string `yaml:"secret"`  // 密钥
}
