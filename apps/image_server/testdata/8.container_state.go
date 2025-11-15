package main

// File: testdata/8.container_state.go
// Description: 获取Docker容器状态的示例代码

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// ListAllContainers 列出所有Docker容器的状态
func ListAllContainers() ([]types.Container, error) {
	// 创建Docker客户端
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("创建Docker客户端失败: %v", err)
	}
	defer cli.Close()

	// 获取所有容器列表（包括停止的容器）
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("获取容器列表失败: %v", err)
	}

	return containers, nil
}

// GetContainerStatus 根据容器名获取容器状态
func GetContainerStatus(containerName string) (types.Container, error) {
	// 创建Docker客户端
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return types.Container{}, fmt.Errorf("创建Docker客户端失败: %v", err)
	}
	defer cli.Close()

	// 使用过滤器按名称查找容器
	filter := filters.NewArgs()
	filter.Add("name", containerName)

	containers, err := cli.ContainerList(ctx, container.ListOptions{
		Filters: filter,
		All:     true,
	})
	if err != nil {
		return types.Container{}, fmt.Errorf("获取容器列表失败: %v", err)
	}

	if len(containers) == 0 {
		return types.Container{}, fmt.Errorf("未找到名为 %s 的容器", containerName)
	}

	// 返回第一个匹配的容器（容器名应该是唯一的）
	return containers[0], nil
}

func main() {
	// 示例：列出所有容器
	allContainers, err := ListAllContainers()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("所有容器状态:")
	for _, container := range allContainers {
		fmt.Printf("ID: %s, 名称: %s, 状态: %s\n", container.ID[:12], container.Names[0][1:], container.State)
	}

	// 示例：获取特定容器状态
	containerName := "hy-dianliadmin" // 替换为实际容器名
	container, err := GetContainerStatus(containerName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n容器 %s 的状态: %s\n", containerName, container.State)
}
