package models

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// File: models/node_model.go
// Description: 定义节点信息的数据模型及其与网络、资源、系统信息的关联关系。

// 节点模型
type NodeModel struct {
	Model
	Title        string         `gorm:"size:64" json:"title"`              // 节点名称
	Uid          string         `gorm:"size:64" json:"uid"`                // 节点UID
	IP           string         `gorm:"size:32" json:"ip"`                 // 节点IP
	Mac          string         `gorm:"size:64" json:"mac"`                // 节点MAC
	Status       int8           `json:"status"`                            // 节点状态
	NetCount     int            `json:"netCount"`                          // 节点网络接口数量
	HoneyIPCount int            `json:"honeyIPCount"`                      // 节点诱捕IP数量
	Resource     NodeResource   `gorm:"serializer:json" json:"resource"`   // 节点资源信息
	SystemInfo   NodeSystemInfo `gorm:"serializer:json" json:"systemInfo"` // 节点系统信息
}

func (n *NodeModel) BeforeDelete(tx *gorm.DB) error {
	// 诱捕转发
	var list []HoneyPortModel
	err := tx.Find(&list, "node_id = ?", n.ID).Delete(&list).Error
	if err != nil {
		return err
	}
	logrus.Infof("关联诱捕转发 %d", len(list))

	// 诱捕ip
	var ipList []HoneyIpModel
	err = tx.Find(&list, "node_id = ?", n.ID).Delete(&ipList).Error
	if err != nil {
		return err
	}
	logrus.Infof("关联诱捕ip %d", len(ipList))

	//节点网络
	var netList []NetModel
	err = tx.Find(&list, "node_id = ?", n.ID).Delete(&netList).Error
	if err != nil {
		return err
	}
	logrus.Infof("关联网络 %d", len(netList))

	// 节点网卡
	var networkList []NodeNetworkModel
	err = tx.Find(&list, "node_id = ?", n.ID).Delete(&networkList).Error
	if err != nil {
		return err
	}
	logrus.Infof("关联节点网卡 %d", len(networkList))
	// 节点
	return nil
}

// 节点资源信息
type NodeResource struct {
	CpuCount              int     `json:"cpuCount"`             // CPU核心数
	CpuUseRate            float64 `json:"cpuUseRate"`           // CPU使用率
	MemTotal              int64   `json:"memTotal"`             // 内存总量
	MemUseRate            float64 `json:"memUseRate"`           // 内存使用率
	DiskTotal             int64   `json:"diskTotal"`            // 磁盘总量
	DiskUseRate           float64 `json:"diskUseRate"`          // 磁盘使用率
	NodePath              string  `json:"nodePath"`             // 节点的部署目录
	NodeResourceOccupancy int64   `json:"nodeResourceOccupany"` // 节点资源磁盘占用
}

// 节点系统信息
type NodeSystemInfo struct {
	HostName            string `json:"hostname"`            // 主机名称
	DistributionVersion string `json:"distributionVersion"` // 发行版本
	CoreVersion         string `json:"coreVersion"`         // 内核版本
	SystemType          string `json:"systemType"`          // 系统类型
	StartTime           string `json:"startTime"`           // 启动时间
	NodeVersion         string `json:"nodeVersion"`         // 节点版本
	NodeCommit          string `json:"nodeCommit"`          // 节点提交
}
