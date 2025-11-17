package main

// File:testdata/3.ip.go
// Description: 测试获取ip地址信息

import (
	"fmt"
	"honey_server/internal/core"
)

func main() {
	core.InitIPDB()
	fmt.Println(core.GetIpAddr("223.104.196.44"))
	fmt.Println(core.GetIpAddr("104.194.71.191"))
}
