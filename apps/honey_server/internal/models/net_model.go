package models

import (
	"errors"
	"fmt"
	"net"

	"gorm.io/gorm"
)

// File: models/net_model.go
// Description: 定义网络配置及扫描状态的数据模型及其与节点、主机的关联关系。

// 网络模型
type NetModel struct {
	Model
	NodeID             uint      `json:"nodeID"`                             // 关联的节点ID
	NodeModel          NodeModel `gorm:"foreignKey:NodeID" json:"-"`         // 关联的节点模型
	Title              string    `gorm:"size:32" json:"title"`               // 网络名称
	Network            string    `gorm:"size:32" json:"network"`             // 网卡名称
	IP                 string    `gorm:"size:32" json:"ip"`                  // 探针IP
	Mask               int8      `json:"mask"`                               // 子网掩码 8-32
	Gateway            string    `gorm:"size:32" json:"gateway"`             // 网关
	HostCount          int       `json:"hostCount"`                          // 子网中活跃的主机数量(存活资产)
	HoneyIpCount       int       `json:"honeyIpCount"`                       // 诱捕IP数量
	ScanStatus         int8      `json:"scanStatus"`                         // 扫描状态
	ScanProgress       float64   `json:"scanProgress"`                       // 扫描进度
	CanUseHoneyIPRange string    `gorm:"size:256" json:"canUseHoneyIPRange"` // 能够使用的诱捕IP范围
}

// 返回当前网络模型的CIDR格式子网表示
func (model NetModel) Subnet() string {
	return fmt.Sprintf("%s/%d", model.IP, model.Mask)
}

// 判断给定的IP地址是否在当前网络模型定义的子网范围内
func (model NetModel) InSubnet(ip string) bool {
	_, _net, _ := net.ParseCIDR(model.Subnet())
	return _net.Contains(net.ParseIP(ip))
}

func (model NetModel) BeforeDelete(tx *gorm.DB) error {
	// 是否有诱捕ip
	var count int64
	tx.Model(HoneyIpModel{}).Where("net_id = ?", model.ID).Count(&count)
	if count > 0 {
		return errors.New("存在诱捕ip，不能删除网络")
	}
	// 将启用的网卡，状态归位
	var nodeNet NodeNetworkModel
	err := tx.Take(&nodeNet, "node_id = ? and network = ?", model.NodeID, model.Network).Error
	if err != nil {
		return nil
	}
	// 修改状态
	tx.Model(&nodeNet).Update("status", 2)
	return nil
}
