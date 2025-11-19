package net_api

// File: api/net_api/use_ip_list.go
// Description: 网络可用ip列表API

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/ip"
	"honey_server/internal/utils/res"
	"net"

	"github.com/gin-gonic/gin"
)

// NetUseIPListResponse IP使用列表响应结构体
type NetUseIPListResponse struct {
	Total              int      `json:"total"`              // IP总数
	Used               int      `json:"used"`               // 已使用IP数
	UseIPList          []string `json:"useIPList"`          // 可用IP列表
	CanUseHoneyIPRange string   `json:"canUseHoneyIPRange"` // 能够使用的诱捕ip范围
}

// NetUseIPListView 查询网络IP使用情况的接口处理函数
func (NetApi) NetUseIPListView(c *gin.Context) {
	// 从请求中绑定并获取ID参数（models.IDRequest结构体）
	cr := middleware.GetBind[models.IDRequest](c)

	var model models.NetModel
	// 根据ID从数据库查询网络信息
	err := global.DB.Take(&model, cr.Id).Error
	if err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}

	// 如果网络的可用诱捕IP范围未初始化，则进行计算
	if model.CanUseHoneyIPRange == "" {
		// 解析子网CIDR格式，获取IP网段信息
		_, ipNet, err := net.ParseCIDR(model.Subnet())
		if err != nil {
			res.FailWithMsg("无效的子网格式", c)
			return
		}

		// 计算可用IP范围（排除网络地址和广播地址）
		startIP := ip.IncrementIP(ipNet.IP)                         // 起始IP（网络地址+1）
		endIP := ip.DecrementIP(ip.BroadcastIP(ipNet))              // 结束IP（广播地址-1）
		model.CanUseHoneyIPRange = ip.FormatIPRange(startIP, endIP) // 格式化IP范围字符串
	}

	// 解析可使用的IP范围为具体的IP列表
	ipList, err := ip.ParseIPRange(model.CanUseHoneyIPRange)
	if err != nil {
		res.FailWithMsg("无效的IP范围", c)
		return
	}

	// 查询已使用的IP：分别从主机表和诱捕IP表获取当前网络下的IP
	var filterIPList1, filterIPList2 []string
	global.DB.Model(models.HostModel{}).Where("net_id = ?", cr.Id).Select("ip").Scan(&filterIPList1)
	global.DB.Model(models.HoneyIpModel{}).Where("net_id = ?", cr.Id).Select("ip").Scan(&filterIPList2)

	// 合并已使用IP列表，用map去重
	usedIPs := make(map[string]struct{})
	for _, ip := range filterIPList1 {
		usedIPs[ip] = struct{}{}
	}
	for _, ip := range filterIPList2 {
		usedIPs[ip] = struct{}{}
	}

	// 生成可用IP列表：过滤掉已使用的IP
	var availableIPs []string
	for _, ip := range ipList {
		if _, exists := usedIPs[ip]; !exists {
			availableIPs = append(availableIPs, ip)
		}
	}

	// 组装响应数据并返回成功结果
	res.OkWithData(NetUseIPListResponse{
		Total:              len(ipList),
		Used:               len(ipList) - len(availableIPs),
		UseIPList:          availableIPs,
		CanUseHoneyIPRange: model.CanUseHoneyIPRange,
	}, c)
}
