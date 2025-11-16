package info

// File: utils/info/network.go
// Description: 获取网卡信息

import (
	"net"
	"strings"
)

// NetworkInfo 网卡信息结构体
type NetworkInfo struct {
	Network string // 网卡名称（如eth0、ens33等）
	Ip      string // 接口的IPv4地址
	Mask    int    // 子网掩码长度
	Net     string // 网络地址（CIDR格式，如192.168.1.0/24）
}

// GetNetworkList 获取网卡信息列表
func GetNetworkList(filterNetworkName string) (list []NetworkInfo, err error) {
	// 获取系统所有网卡
	faces, err := net.Interfaces()
	if err != nil {
		return
	}

	// 遍历每个网卡
	for _, face := range faces {
		faceName := face.Name

		// 跳过回环接口（lo）
		if faceName == "lo" {
			continue
		}

		// 过滤掉名称以指定前缀的网卡（如诱捕相关网卡）
		if strings.HasPrefix(faceName, filterNetworkName) {
			continue
		}

		// 获取当前接口的所有地址
		addrs, err := face.Addrs()
		if err != nil {
			continue // 获取地址失败则跳过当前接口
		}

		// 遍历接口的每个地址
		for _, addr := range addrs {
			// 解析地址为CIDR格式（如192.168.1.100/24）
			ip, _net, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue // 解析失败则跳过当前地址
			}

			// 只保留IPv4地址（过滤IPv6）
			if ip.To4() == nil {
				continue
			}

			// 获取子网掩码长度
			mask, _ := _net.Mask.Size()

			// 将当前网卡信息添加到结果列表
			list = append(list, NetworkInfo{
				Network: faceName,
				Ip:      ip.String(),
				Mask:    mask,
				Net:     _net.String(),
			})
		}
	}
	return
}
