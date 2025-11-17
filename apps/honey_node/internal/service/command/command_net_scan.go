package command

// File: service/command/net_scan.go
// Description: 节点客户端中处理网络扫描命令的逻辑实现，负责接收网络扫描请求并返回扫描进度及结果

import (
	"fmt"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/ip"
	"net"
	"sync"
	"time"

	"github.com/j-keck/arping"
)

// CmdNetScan 处理网络扫描命令
// 解析扫描IP范围，通过ARP协议并发扫描指定IP，实时上报扫描进度、IP与MAC对应关系，完成后发送结束标识
func (nc *NodeClient) CmdNetScan(request *node_rpc.CmdRequest) {
	// 提取网络扫描请求的具体参数（IP范围、过滤IP、网络接口等）
	req := request.GetNetScanInMessage()
	fmt.Printf("网络扫描 %v\n", req)
	startTime := time.Now() // 记录扫描开始时间

	// 解析IP范围字符串，转换为具体IP列表
	ipList, err := ip.ParseIPRange(req.IpRange)
	if err != nil {
		// 解析失败时，发送包含错误信息的结束响应
		fmt.Println(err)
		nc.cmdResponseChan <- &node_rpc.CmdResponse{
			CmdType: node_rpc.CmdType_cmdNetScanType,
			TaskID:  request.TaskID,
			NodeID:  nc.config.System.Uid,
			NetScanOutMessage: &node_rpc.NetScanOutMessage{
				End:      true,                              // 标识扫描结束
				Progress: 0,                                 // 进度0%
				NetID:    req.NetID,                         // 关联的网络ID
				ErrMsg:   fmt.Sprintf("解析扫描ip列表出错 %s", err), // 错误信息
			},
		}
		return
	}

	// 将过滤IP列表转换为map，便于快速判断是否需要跳过该IP
	filterIPList := map[string]struct{}{}
	for _, s := range req.FilterIPList {
		filterIPList[s] = struct{}{}
	}

	// 扫描配置：指定网络接口、最大并发数
	iface := req.Network          // 扫描使用的网络接口名称
	concurrency := 200            // 最大并发扫描数，控制资源占用
	totalIPs := len(ipList)       // 总IP数量，用于计算进度
	processed := 0                // 已处理的IP数量
	var processedMutex sync.Mutex // 保护processed变量的互斥锁，确保并发安全

	// 创建信号量通道，控制并发数量（最多同时运行concurrency个goroutine）
	semaphore := make(chan struct{}, concurrency)

	fmt.Printf("开始扫描 %d 个IP地址，并发数: %d\n", totalIPs, concurrency)

	var wg sync.WaitGroup // 等待组，用于等待所有扫描goroutine完成

	// 遍历IP列表，对每个IP启动goroutine进行扫描
	for _, ipStr := range ipList {
		// 跳过过滤列表中的IP
		if _, exists := filterIPList[ipStr]; exists {
			continue
		}

		wg.Add(1) // 增加等待组计数

		// 获取信号量（若通道已满则阻塞，限制并发数）
		semaphore <- struct{}{}

		// 启动goroutine执行ARP扫描
		go func(ip string) {
			defer wg.Done() // 扫描完成后减少等待组计数
			defer func() {
				<-semaphore // 释放信号量，允许新的goroutine执行
			}()

			// 对目标IP执行ARP扫描，获取MAC地址
			mac, _, err := arping.PingOverIfaceByName(net.ParseIP(ip), iface)

			// 更新已处理数量和进度（加锁保证线程安全）
			processedMutex.Lock()
			processed++
			progress := float64(processed) / float64(totalIPs) * 100 // 计算进度百分比
			processedMutex.Unlock()

			// 若扫描出错（如目标IP无响应），直接返回
			if err != nil {
				return
			}

			// 打印扫描结果
			fmt.Printf("%s %s %.2f\n", ip, mac, progress)

			// 发送中间结果响应：包含当前IP、MAC、进度
			nc.cmdResponseChan <- &node_rpc.CmdResponse{
				CmdType: node_rpc.CmdType_cmdNetScanType,
				TaskID:  request.TaskID,
				NodeID:  nc.config.System.Uid,
				NetScanOutMessage: &node_rpc.NetScanOutMessage{
					End:      false,             // 未结束
					Progress: float32(progress), // 当前进度
					NetID:    req.NetID,         // 关联网络ID
					Ip:       ip,                // 扫描到的IP
					Mac:      mac.String(),      // 对应的MAC地址
				},
			}
		}(ipStr)
	}

	wg.Wait()        // 等待所有扫描goroutine完成
	close(semaphore) // 关闭信号量通道

	// 发送扫描完成响应
	nc.cmdResponseChan <- &node_rpc.CmdResponse{
		CmdType: node_rpc.CmdType_cmdNetScanType,
		TaskID:  request.TaskID,
		NodeID:  nc.config.System.Uid,
		NetScanOutMessage: &node_rpc.NetScanOutMessage{
			End:      true,      // 标识扫描结束
			Progress: 100,       // 进度100%
			NetID:    req.NetID, // 关联网络ID
		},
	}

	// 打印扫描总耗时
	fmt.Printf("\n扫描完成，耗时: %v\n", time.Since(startTime))
}
