package global

// File: global/enter.go
// Description: 声明全局变量，供其他模块共享。

import (
	"image_server/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Version   = "v1.0.1"              // 版本号
	Commit    = "56a9c63f"            // 提交ID
	BuildTime = "2025-11-10 17:12:34" // 构建时间
)

var (
	DB     *gorm.DB       // 数据库实例
	Config *config.Config // 配置实例
	Log    *logrus.Entry  // 日志实例
	Redis  *redis.Client  // Redis 客户端实例
)
