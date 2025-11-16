package config

// File: config/enter.go
// Description: 定义诱捕服务所需的配置和相关方法。

// 统一管理配置结构体
type Config struct {
	Logger Logger `yaml:"logger"`
	System System `yaml:"system"`
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
