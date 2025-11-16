package main

// File: testdata/1.ip.go
// Description: 获取网卡信息

import (
	"fmt"
	"honey_node/internal/utils/ip"
)

func main() {
	ifaceName := "eth0"
	fmt.Println(ip.GetNetworkInfo(ifaceName))
}
