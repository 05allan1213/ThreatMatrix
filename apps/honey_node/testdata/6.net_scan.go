package main

// File: testdata/6.net_scan.go
// Description: 子网扫描

import (
	"fmt"
	"honey_node/internal/utils/ip"
	"net"
	"sync"
	"time"

	"github.com/j-keck/arping"
)

func main() {
	t1 := time.Now()
	fmt.Println(time.Now().Format(time.DateTime))
	ipList, err := ip.ParseIPRange("10.2.0.1-10.2.0.5")
	if err != nil {
		fmt.Println(err)
		return
	}
	iface := "eth0"
	wait := sync.WaitGroup{}
	for _, s := range ipList {
		wait.Add(1)
		go func(s string) {
			defer wait.Done()
			mac, _, err := arping.PingOverIfaceByName(net.ParseIP(s), iface)
			if err != nil {
				fmt.Println(s, err)
				return
			}
			fmt.Println(s, mac)
		}(s)
	}
	wait.Wait()
	fmt.Println("扫描完成", time.Since(t1))
	fmt.Println(time.Now().Format(time.DateTime))
}