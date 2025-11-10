package global

// File: global/enter.go
// Description: 声明全局变量，供其他模块共享。

import (
	"honey_server/config"

	"gorm.io/gorm"
)

var (
	DB     *gorm.DB       // 数据库实例
	Config *config.Config // 配置实例
)
