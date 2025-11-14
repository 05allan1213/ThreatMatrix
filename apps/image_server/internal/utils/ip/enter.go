package ip

// File:utils/ip/enter.go
// Description: 判断ip是否是本地ip

import "net"

// 判断ip是否是本地ip
func HasLocalIPAddr(_ip string) bool {
	ip := net.ParseIP(_ip)
	// 判断是否是私有地址
	if ip.IsPrivate() {
		return true
	}
	// 判断是否是回环地址
	if ip.IsLoopback() {
		return true
	}
	return false
}
