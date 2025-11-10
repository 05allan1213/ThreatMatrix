package config

// File: config/enter.go
// Description: 定义诱捕服务所需的配置和相关方法。

import "fmt"

// 配置结构体
type Config struct {
	DB     DB     `yaml:"db"`
	Logger Logger `yaml:"logger"`
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

// 日志配置
type Logger struct {
	Format  string `yaml:"format"`
	Level   string `yaml:"level"`
	AppName string `yaml:"appName"`
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
