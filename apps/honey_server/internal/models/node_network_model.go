package models

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// File: models/node_network_model.go
// Description: 定义节点网卡信息的数据模型及其与节点的关联关系。

// 节点网卡模型
type NodeNetworkModel struct {
	Model
	NodeID    uint      `json:"nodeID"`                     // 关联的节点ID
	NodeModel NodeModel `gorm:"foreignKey:NodeID" json:"-"` // 关联的节点模型
	Network   string    `gorm:"size:32" json:"network"`     // 网卡名称
	IP        string    `gorm:"size:32" json:"ip"`          // 探针IP
	Mask      int8      `gorm:"size:32" json:"mask"`        // 子网掩码 8-32
	Gateway   string    `gorm:"size:32" json:"gateway"`     // 网关
	Status    int8      `json:"status"`                     // 网关状态 1 启用 2 未启用
}

func (n *NodeNetworkModel) BeforeDelete(tx *gorm.DB) error {
	// 先找有没有网络
	if n.Status == 2 {
		// 未启用
		return nil
	}
	var net NetModel
	err := tx.Take(&net, "node_id = ? and network = ?", n.NodeID, n.Network).Error
	if err != nil {
		// 未启用
		return nil
	}

	// 判断有没有诱捕ip
	var count int64
	tx.Model(HoneyIpModel{}).Where("net_id = ?", net.ID).Count(&count)
	if count > 0 {
		return errors.New("此网卡的网络存在诱捕ip，不可删除")
	}

	// 关联删除网络表和主机表
	var hostList []HostModel
	tx.Find(&hostList, "net_id = ?", net.ID).Delete(&hostList)
	tx.Delete(&net)
	logrus.Infof("关联删除主机记录 %d", len(hostList))
	logrus.Infof("关联删除网络 %s", net.Title)
	return nil
}
