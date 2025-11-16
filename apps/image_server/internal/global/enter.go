package global

// File: global/enter.go
// Description: 声明全局变量，供其他模块共享。

import (
	"image_server/internal/config"

	"github.com/docker/docker/client"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	Version   = "v1.0.1"             // 版本号
	Commit    = "5dd146421917"       // 提交ID
	BuildTime = "2025-11-16 9:45:25" // 构建时间
)

var (
	DB           *gorm.DB       // 数据库实例
	Config       *config.Config // 配置实例
	Log          *logrus.Entry  // 日志实例
	Redis        *redis.Client  // Redis 客户端实例
	DockerClient *client.Client // Docker 客户端实例
)
