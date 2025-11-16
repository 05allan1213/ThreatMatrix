package info

// File:utils/info/enter.go
// Description: 系统资源信息获取

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

// 系统资源信息结构体
type ResourceMessage struct {
	CpuCount              int64   `json:"cpuCount,omitempty"`              // CPU核心数
	CpuUseRate            float32 `json:"cpuUseRate,omitempty"`            // CPU使用率
	MemTotal              int64   `json:"memTotal,omitempty"`              // 总内存大小
	MemUseRate            float32 `json:"memUseRate,omitempty"`            // 内存使用率
	DiskTotal             int64   `json:"diskTotal,omitempty"`             // 磁盘总容量
	DiskUseRate           float32 `json:"diskUseRate,omitempty"`           // 磁盘使用率
	NodePath              string  `json:"nodePath,omitempty"`              // 节点的部署目录
	NodeResourceOccupancy int64   `json:"nodeResourceOccupancy,omitempty"` // 节点磁盘占用
}

// GetResourceInfo 获取系统资源信息
func GetResourceInfo(nodePath string) (*ResourceMessage, error) {
	// 获取CPU核心数（true表示包含逻辑核心）
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	// 获取CPU使用率（采样时间1秒，false表示返回整体使用率而非每个核心）
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	// 获取虚拟内存信息（总内存、使用率等）
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// 获取节点路径（nodePath）的磁盘使用信息（用于计算该路径的占用空间）
	nodeDiskInfo, err := disk.Usage(nodePath)
	if err != nil {
		return nil, err
	}

	// 获取根目录（"/"）的磁盘信息（用于获取总磁盘容量和整体使用率）
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	// 组装资源信息结构体并返回
	message := &ResourceMessage{
		CpuCount:              int64(cpuCount),
		CpuUseRate:            float32(cpuPercent[0]), // cpu.Percent返回切片，取第一个元素（整体使用率）
		MemTotal:              int64(memInfo.Total),
		MemUseRate:            float32(memInfo.UsedPercent),
		DiskTotal:             int64(diskInfo.Total),
		DiskUseRate:           float32(diskInfo.UsedPercent),
		NodePath:              nodePath,
		NodeResourceOccupancy: int64(nodeDiskInfo.Used),
	}

	return message, nil
}
