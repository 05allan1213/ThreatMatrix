package main

// File: testdata/1.net_segment.go
// Description: 测试网段

import (
	"fmt"
	"honey_server/internal/utils/ip"
)

func main() {
	fmt.Println(ip.ParseCIDRGetUseIPRange("192.168.100.1/24"))
	fmt.Println(ip.ParseCIDRGetUseIPRange("192.168.100.1/25"))
	fmt.Println(ip.ParseCIDRGetUseIPRange("192.168.100.1/16"))
}
