package main

import (
	"fmt"
	"honey_server/internal/utils/ip"
)

func main() {
	ipList, err := ip.ParseIPRange("192.168.200.2-192.168.200.3,192.168.200.5,192.168.200.240")
	fmt.Println(ipList, err)
}
