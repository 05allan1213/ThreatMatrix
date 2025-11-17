package main

// File: testdata/5.arp.go
// Description: 测试arping

import (
	"fmt"
	"net"

	"github.com/j-keck/arping"
)

func main() {
	mac, t, err := arping.Ping(net.ParseIP("10.2.0.1"))
	fmt.Println(mac, t, err)
	mac, t, err = arping.PingOverIfaceByName(net.ParseIP("10.2.0.1"), "eth0")
	fmt.Println(mac, t, err)
}