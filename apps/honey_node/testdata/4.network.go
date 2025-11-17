package main

// File: testdata/4.network.go
// Description: 获取指定网卡列表

import (
	"fmt"
	"honey_node/internal/utils/info"
)

func main() {
	networkList, err := info.GetNetworkList([]string{"br-", "vir", "mc_", "docker"})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, networkInfo := range networkList {
		fmt.Println(networkInfo)
	}
}
