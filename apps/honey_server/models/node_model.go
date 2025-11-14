package models

// File: models/node_model.go
// Description: 定义节点信息的数据模型及其与网络、资源、系统信息的关联关系。

// 节点模型
type NodeModel struct {
	Model
	Title          string         `gorm:"size:32" json:"title"`                  // 节点名称
	IP             string         `gorm:"size:32" json:"ip"`                     // 节点IP
	Status         int8           `json:"status"`                                // 节点状态
	NetCount       int            `json:"netCount"`                              // 节点网络接口数量
	HoneyIPCount   int            `json:"honeyIPCount"`                          // 节点诱捕IP数量
	Resource       NodeResource   `gorm:"serializer:json" json:"resource"`       // 节点资源信息
	NodeSystemInfo NodeSystemInfo `gorm:"serializer:json" json:"nodeSystemInfo"` // 节点系统信息
}

// 节点资源信息
type NodeResource struct {
	CpuCount           int     `json:"cpuCount"`           // CPU核心数
	CpuUseRate         float64 `json:"cpuUseRate"`         // CPU使用率
	MemTotal           int64   `json:"memTotal"`           // 内存总量
	MemUseRate         float64 `json:"memUseRate"`         // 内存使用率
	DiskTotal          int64   `json:"diskTotal"`          // 磁盘总量
	DiskUseRate        float64 `json:"diskUseRate"`        // 磁盘使用率
	NodePath           string  `json:"nodePath"`           // 节点的部署目录
	NodeResourceOccupy int64   `json:"nodeResourceOccupy"` // 节点目录的资源占用
}

// 节点系统信息
type NodeSystemInfo struct {
	Hostname            string `json:"hostname"`            // 主机名称
	DistributionVersion string `json:"distributionVersion"` // 发行版本
	CoreVersion         string `json:"coreVersion"`         // 内核版本
	SystemType          string `json:"systemType"`          // 系统类型
	StartTime           string `json:"startTime"`           // 启动时间
}
