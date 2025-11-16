package ip

// File:utils/ip/enter.go
// Description: 获取指定网卡的IPv4地址和MAC地址

import (
	"fmt"
	"net"
)

// 根据网卡名称获取该网卡的IPv4地址和MAC地址
func GetNetworkInfo(i string) (ip string, mac string, err error) {
	// 根据网卡名称获取网络接口信息
	iface, err := net.InterfaceByName(i)
	if err != nil {
		err = fmt.Errorf("无法获取网卡 %s: %v", i, err)
		return
	}

	// 获取该网卡的所有网络地址
	addrs, err := iface.Addrs()
	if err != nil {
		err = fmt.Errorf("无法获取网卡 %s 的地址: %s", iface.Name, err)
		return
	}

	// 提取网卡的MAC地址
	mac = iface.HardwareAddr.String()

	// 遍历地址列表，筛选出IPv4地址
	for _, addr := range addrs {
		var _ip net.IP
		// 从地址中提取IP部分（支持IPNet和IPAddr两种类型）
		switch v := addr.(type) {
		case *net.IPNet:
			_ip = v.IP
		case *net.IPAddr:
			_ip = v.IP
		}
		// 检查是否为IPv4地址（To4()返回非nil表示是IPv4）
		if _ip.To4() != nil {
			ip = _ip.String()
			// 若存在多个IPv4地址，此处会覆盖前一个，最终返回最后一个找到的IPv4地址
		}
	}

	// 若未找到IPv4地址，返回错误
	if ip == "" {
		err = fmt.Errorf("%s 此接口无IPv4地址", iface.Name)
		return
	}

	return
}
