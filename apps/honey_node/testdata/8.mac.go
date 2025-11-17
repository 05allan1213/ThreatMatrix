package main

// File:testdata/8.mac.go
// Description: 获取MAC地址的厂商信息

import (
	"fmt"
	"honey_node/internal/core"
)

func main() {
	fmt.Println(core.ManufQuery("02:42:09:93:de:f7"))
	fmt.Println(core.ManufQuery("02:42:1a:19:4e:fb"))
	fmt.Println(core.ManufQuery("02:42:eb:6d:b0:2a"))
	fmt.Println(core.ManufQuery("02:42:43:c3:5c:aa"))
	fmt.Println(core.ManufQuery("52:54:00:c2:18:ec"))
	fmt.Println(core.ManufQuery("02:41:4e:e2:54:30"))
	fmt.Println(core.ManufQuery("92:c3:36:e5:ee:48"))
	fmt.Println(core.ManufQuery("00:50:56:c0:00:08"))
	fmt.Println(core.ManufQuery("00:50:56:e4:1f:23"))
	fmt.Println(core.ManufQuery("30-9C-23-43-7E-61"))
	fmt.Println(core.ManufQuery("30-1C-23-43-7E-61"))
}
