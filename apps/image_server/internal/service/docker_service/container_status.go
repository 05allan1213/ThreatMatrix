package docker_service

// File: service/docker_service/container_status.go
// Description: 提供 Docker 容器的查询服务，包括列出所有容器与根据名称查询容器状态。

import (
	"context"
	"fmt"
	"image_server/internal/global"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// ListAllContainers 列出所有 Docker 容器（包括已停止的）
func ListAllContainers() ([]container.Summary, error) {
	// 调用 Docker 客户端获取容器列表
	containers, err := global.DockerClient.ContainerList(context.Background(), container.ListOptions{
		All: true, // true 表示包含已停止的容器
	})
	if err != nil {
		return nil, fmt.Errorf("获取容器列表失败: %v", err)
	}

	return containers, nil
}

// PrefixContainerStatus 根据容器名前缀获取容器状态
func PrefixContainerStatus(containerName string) (summaryList []container.Summary, err error) {
	// 创建过滤器：按容器名称过滤
	filter := filters.NewArgs()
	filter.Add("name", containerName)

	// 查询匹配的容器
	containers, err := global.DockerClient.ContainerList(context.Background(), container.ListOptions{
		Filters: filter,
		All:     true, // 包含停止容器
	})
	if err != nil {
		return
	}

	// 返回首个匹配项（正常情况下名称唯一）
	return containers, nil
}
