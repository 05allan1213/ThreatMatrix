package core

// File: core/redis.go
// Description: 实现Redis初始化，用于建立与Redis服务器的连接。

import (
	"context"
	"honey_server/internal/global"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var client *redis.Client

// 初始化Redis连接
func InitRedis() (client *redis.Client) {
	conf := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})
	// 测试连接
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatalf("连接redis失败 %s", err)
		return
	}
	logrus.Infof("成功连接redis")
	return rdb
}

var onceRedis sync.Once

// 获取Redis客户端实例（单例模式）
func GetRedisClient() *redis.Client {
	onceRedis.Do(func() {
		client = InitRedis()
	})
	return client
}
