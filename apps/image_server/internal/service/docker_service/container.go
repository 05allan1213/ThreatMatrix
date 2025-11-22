package docker_service

// File: service/docker_service/container.go
// Description: 使用 Docker SDK 创建并启动容器，支持指定网络、静态 IP 与镜像参数。

import (
	"context"
	"image_server/internal/global"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// RunContainer 创建并启动容器，返回容器 ID
func RunContainer(containerName, networkName, ip, image string) (containerID string, err error) {

	// 容器基本配置：指定镜像
	containerConfig := &container.Config{
		Image: image,
	}

	// 主机配置：设置网络模式，不自动删除容器
	hostConfig := &container.HostConfig{
		AutoRemove:  false,
		NetworkMode: container.NetworkMode(networkName),
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
	}

	// 网络配置：指定 IP 和网络
	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {
				IPAMConfig: &network.EndpointIPAMConfig{
					IPv4Address: ip,
				},
			},
		},
	}

	// 创建容器
	createResp, err := global.DockerClient.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		networkingConfig,
		nil,
		containerName,
	)
	if err != nil {
		return
	}

	// 启动容器
	err = global.DockerClient.ContainerStart(context.Background(), createResp.ID, container.StartOptions{})
	if err != nil {
		return
	}

	// 返回 12 位容器 ID
	containerID = createResp.ID[:12]
	return
}
