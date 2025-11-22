package ip_service

// File: service/ip_service/ip_load.go
// Description: 负责从数据库加载IP配置，确保IP网络接口的持久化

import (
	"honey_node/internal/global"
	"honey_node/internal/models"
	"honey_node/internal/utils"
	"honey_node/internal/utils/info"

	"github.com/sirupsen/logrus"
)

// IPLoad 加载数据库中的IP配置并与系统实际状态同步
func IPLoad() {
	// 从数据库加载所有持久化的IP配置记录
	var ipList []models.IpModel
	global.DB.Find(&ipList)

	// 获取当前系统的实际网络接口信息（接口名称→IP列表的映射）
	networkMap, err := info.GetNetworkInterfaces()
	if err != nil {
		logrus.Fatalf("获取网卡错误 %s", err)
		return
	}

	// 遍历每条IP配置记录，与系统实际状态对比并同步
	for _, model := range ipList {
		// 检查配置的网络接口是否存在于当前系统中
		ips, ok := networkMap[model.LinkName]
		if ok {
			// 接口存在时，校验配置的IP是否与系统实际IP一致
			if !utils.InList(ips, model.Ip) {
				logrus.Errorf("网卡 %s 对应的ip地址错误 %v %s", model.LinkName, ips, model.Ip)
				continue // IP不一致时跳过重建（或可扩展为自动修复）
			}
			continue // 接口存在且IP正确，无需处理
		}

		// 接口不存在时，根据数据库配置重建网络接口和IP
		_, err := SetIp(SetIpRequest{
			Ip:       model.Ip,       // 配置的IP地址
			Mask:     model.Mask,     // 子网掩码位数
			LinkName: model.LinkName, // 网络接口名称
			Network:  model.Network,  // 基于的主网卡名称
			Mac:      model.Mac,      // 接口MAC地址
		})
		if err != nil {
			logrus.Errorf("初始化ip错误 %s", err)
			continue // 单条配置重建失败不影响其他配置处理
		}
	}
}
