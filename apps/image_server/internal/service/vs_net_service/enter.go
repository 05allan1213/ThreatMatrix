package vs_net_service

// File: service/vs_net_service/enter.go
// Description: 虚拟网络服务，负责初始化虚拟网络配置，检查Docker中是否存在指定网络，若不存在则创建；若存在则验证子网配置是否匹配。

import (
	"context"
	"image_server/internal/global"

	"github.com/docker/docker/api/types/network"
	"github.com/sirupsen/logrus"
)

// 初始化虚拟网络服务
//
// 1. 检查并设置虚拟网络默认配置（名称、子网、容器前缀）
// 2. 查询Docker中所有网络，检查配置的网络是否存在
// 3. 若网络不存在，根据配置创建Docker网络
// 4. 若网络存在，验证其实际子网是否与配置一致，不一致则报错提示
func Run() {
	// 获取全局虚拟网络配置，若配置项为空则设置默认值
	cfg := global.Config.VsNet
	if cfg.Name == "" {
		cfg.Name = "honey-hy" // 默认网络名称
	}
	if cfg.Net == "" {
		cfg.Net = "10.2.0.0/24" // 默认子网配置
	}
	if cfg.Prefix == "" {
		cfg.Prefix = "hy-" // 默认容器名称前缀
	}

	// 调用Docker API获取所有网络列表，用于后续检查
	networks, err := global.DockerClient.NetworkList(context.Background(), network.ListOptions{})
	if err != nil {
		logrus.Fatalf("获取虚拟网络列表失败: %v", err)
	}

	// 遍历网络列表，查找配置中指定名称的网络是否存在
	var found bool
	var existingNetwork network.Summary
	for _, network := range networks {
		if network.Name == cfg.Name {
			found = true
			existingNetwork = network
			break
		}
	}

	// 若网络不存在，根据配置创建Docker网络
	if !found {
		// 配置网络IPAM（IP地址管理）信息，指定子网
		ipam := network.IPAM{
			Driver: "default",
			Config: []network.IPAMConfig{
				{
					Subnet: cfg.Net,
				},
			},
		}
		// 调用Docker API创建网络
		_, err := global.DockerClient.NetworkCreate(context.Background(), cfg.Name, network.CreateOptions{
			IPAM: &ipam,
		})
		if err != nil {
			logrus.Fatalf("创建网络失败: %v", err)
		}
		logrus.Printf("成功创建网络 %s，子网为 %s", cfg.Name, cfg.Net)
		return
	}

	// 网络已存在，检查实际子网是否与配置一致（简化逻辑：假设仅一个IPAM配置）
	if len(existingNetwork.IPAM.Config) > 0 && existingNetwork.IPAM.Config[0].Subnet != cfg.Net {
		logrus.Warnf("警告: 网络 %s 存在，但子网不匹配。现有子网: %s，配置子网: %s",
			cfg.Name, existingNetwork.IPAM.Config[0].Subnet, cfg.Net)
		logrus.Fatalf("请排查网络配置问题")
		return
	}

	// 网络存在且子网匹配，输出成功信息
	logrus.Infof("网络 %s 存在且子网匹配: %s", cfg.Name, cfg.Net)
}
