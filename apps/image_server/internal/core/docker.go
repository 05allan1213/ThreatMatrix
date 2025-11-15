package core

// File: core/docker.go
// Description: 初始化 Docker 客户端

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

// 初始化 Docker 客户端
func InitDocker() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.Fatalf("创建Docker客户端失败: %v", err)
	}
	return cli
}
