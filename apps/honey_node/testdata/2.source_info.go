package main

// File: testdata/2.source_info.go
// Description: 获取节点资源信息

import (
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// ResourceMessage 定义资源信息结构体
type ResourceMessage struct {
	CpuCount              int64   `json:"cpuCount,omitempty"`
	CpuUseRate            float32 `json:"cpuUseRate,omitempty"`
	MemTotal              int64   `json:"memTotal,omitempty"`
	MemUseRate            float32 `json:"memUseRate,omitempty"`
	DiskTotal             int64   `json:"diskTotal,omitempty"`
	DiskUseRate           float32 `json:"diskUseRate,omitempty"`
	NodePath              string  `json:"nodePath,omitempty"`
	NodeResourceOccupancy int64   `json:"nodeResourceOccupancy,omitempty"`
}

// GetResourceInfo 通用函数，用于获取资源信息
func GetResourceInfo(nodePath string) (*ResourceMessage, error) {
	// 获取 CPU 核心数
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	// 获取 CPU 使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 获取磁盘信息
	nodeDiskInfo, err := disk.Usage(nodePath)
	if err != nil {
		return nil, fmt.Errorf("无法访问节点路径 '%s': %w", nodePath, err)
	}

	// 获取磁盘信息
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	// 填充 ResourceMessage 结构体
	message := &ResourceMessage{
		CpuCount:              int64(cpuCount),
		CpuUseRate:            float32(cpuPercent[0]),
		MemTotal:              int64(memInfo.Total),
		MemUseRate:            float32(memInfo.UsedPercent),
		DiskTotal:             int64(diskInfo.Total),
		DiskUseRate:           float32(diskInfo.UsedPercent),
		NodePath:              nodePath,
		NodeResourceOccupancy: int64(nodeDiskInfo.Used),
	}

	return message, nil
}

func main() {
	// 使用当前工作目录作为节点路径
	nodePath, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取当前工作目录失败: %v\n", err)
		return
	}
	
	message, err := GetResourceInfo(nodePath)
	if err != nil {
		fmt.Printf("获取资源信息失败: %v\n", err)
		return
	}

	fmt.Printf("CPU 核心数: %d\n", message.CpuCount)
	fmt.Printf("CPU 使用率: %.2f%%\n", message.CpuUseRate)
	fmt.Printf("内存总量: %d bytes\n", message.MemTotal)
	fmt.Printf("内存使用率: %.2f%%\n", message.MemUseRate)
	fmt.Printf("磁盘总量: %d bytes\n", message.DiskTotal)
	fmt.Printf("磁盘使用率: %.2f%%\n", message.DiskUseRate)
	fmt.Printf("节点路径: %s\n", message.NodePath)
	fmt.Printf("节点资源占用率: %d bytes\n", message.NodeResourceOccupancy)
}