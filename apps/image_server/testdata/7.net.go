package main

// File: testdata/7.net.go
// Description: 计算子网段，检查 IP 是否在某个网段内

import (
	"fmt"
	"net"
)

func main() {
	x := "10.2.0.0/16"                 // 一个 CIDR 网段，掩码长度为 /16，可覆盖 10.2.0.0 ~ 10.2.255.255
	ip, ipNet, err := net.ParseCIDR(x) // 解析 CIDR，返回网络地址和 *IPNet
	fmt.Println(ip, ipNet, err)        // ip=网络地址，ipNet=包含掩码信息

	ip4 := ip.To4() // 转为 IPv4（4 字节）
	ip4[3] = 2      // 手动修改最后一段，将 10.2.0.0 改为 10.2.0.2
	fmt.Println(ip4.String())

	// 检查指定 IP 是否在该 CIDR 网段内
	fmt.Println(ipNet.Contains(net.ParseIP("10.2.1.3"))) // 由于 /16，10.2.1.3 属于该网段 → true
}
