package main

// File: testdata/7.net_scan-well.go
// Description: 子网扫描（优化版）

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
	ipList, err := ip.ParseIPRange("10.2.0.1-10.2.0.255")
	if err != nil {
		fmt.Println(err)
		return
	}

	filterIPList := map[string]struct{}{
		"10.2.0.3": {},
	}

	iface := "eth0"
	concurrency := 100 // 最大并发数
	totalIPs := len(ipList)
	processed := 0
	var processedMutex sync.Mutex

	// 创建并发控制通道
	semaphore := make(chan struct{}, concurrency)

	fmt.Printf("开始扫描 %d 个IP地址，并发数: %d\n", totalIPs, concurrency)

	var wg sync.WaitGroup
	for _, s := range ipList {
		if _, exists := filterIPList[s]; exists {
			continue
		}

		wg.Add(1)

		// 获取信号量（如果已满则阻塞）
		semaphore <- struct{}{}

		go func(s string) {
			defer wg.Done()
			defer func() {
				// 释放信号量
				<-semaphore
			}()

			mac, _, err := arping.PingOverIfaceByName(net.ParseIP(s), iface)

			// 更新进度
			processedMutex.Lock()
			processed++
			progress := float64(processed) / float64(totalIPs) * 100
			fmt.Printf("进度: %.2f%% (%d/%d)\n", progress, processed, totalIPs)
			processedMutex.Unlock()

			if err != nil {
				return
			}

			fmt.Printf("%s 在线，MAC: %s\n", s, mac)
		}(s)
	}

	wg.Wait()
	close(semaphore) // 扫描完成后关闭通道
	fmt.Printf("\n扫描完成，耗时: %v\n", time.Since(t1))
}
