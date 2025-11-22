package info

// File: utils/info/network_map.go
// Description: 提供获取主机网络接口及其IPv4地址的功能，过滤非活跃接口并结构化返回网卡名称与对应IP的映射关系

import (
	"fmt"
	"net"
)

// GetNetworkInterfaces 获取主机所有活跃网络接口及其关联的IPv4地址
func GetNetworkInterfaces() (map[string][]string, error) {
	// 创建映射存储网卡名称与对应IPv4地址列表的关联关系
	interfacesMap := make(map[string][]string)

	// 获取主机上所有的网络接口（包括物理网卡、虚拟网卡、回环接口等）
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网络接口失败: %v", err)
	}

	// 遍历每个网络接口，提取有效IPv4地址信息
	for _, iface := range interfaces {
		// 过滤掉状态为非活跃（down）的接口，只处理up状态的网络接口
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 获取当前接口绑定的所有网络地址（包括IPv4、IPv6等）
		addresses, err := iface.Addrs()
		if err != nil {
			fmt.Printf("获取接口 %s 的地址失败: %v\n", iface.Name, err)
			continue // 单个接口获取失败不影响整体流程，继续处理下一个接口
		}

		// 提取当前接口的IPv4地址
		var ipv4Addresses []string
		for _, addr := range addresses {
			var ip net.IP

			// 处理不同类型的网络地址结构（IPNet包含子网掩码，IPAddr仅为IP地址）
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP // 从IPNet中提取IP地址
			case *net.IPAddr:
				ip = v.IP // 从IPAddr中提取IP地址
			}

			// 过滤出IPv4地址（排除IPv6地址）
			if ip != nil && ip.To4() != nil {
				ipv4Addresses = append(ipv4Addresses, ip.String())
			}
		}

		// 仅当接口存在IPv4地址时，才将其加入结果映射（避免空地址条目）
		if len(ipv4Addresses) > 0 {
			interfacesMap[iface.Name] = ipv4Addresses
		}
	}

	return interfacesMap, nil
}
