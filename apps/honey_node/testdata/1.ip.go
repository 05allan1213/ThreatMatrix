package main

import (
	"fmt"
	"honey_node/internal/utils/ip"
)

func main() {
	ifaceName := "eth0"
	fmt.Println(ip.GetNetworkInfo(ifaceName))
}
