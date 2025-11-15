package main

// File: main.go
// Description: 使用 Docker SDK 创建并启动容器，包含指定网络与固定 IP 配置。

import (
	"fmt"
	"image_server/internal/core"
	"image_server/internal/global"
	"image_server/internal/service/docker_service"
)

func main() {
	global.DockerClient = core.InitDocker()

	// 输入参数
	containerName := "alpine" // 容器名称
	image := "alpine:latest"  // 镜像名称
	networkName := "honey-hy" // 目标网络
	ipAddress := "10.2.0.2"   // 静态 IP

	id, err := docker_service.RunContainer(containerName, networkName, ipAddress, image)
	fmt.Println(id, err)
}
