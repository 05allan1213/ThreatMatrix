package net_api

// File: net_api.go
// Description: 网络相关API接口处理逻辑，包含网络信息更新等操作的实现

import (
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/ip"
	"honey_server/internal/utils/res"
	"net"

	"github.com/gin-gonic/gin"
)

// UpdateRequest 网络信息更新的请求参数
type UpdateRequest struct {
	ID                 uint   `json:"id" binding:"required"`    // 网络ID，必填，用于指定要更新的网络
	Title              string `json:"title" binding:"required"` // 网络标题，必填
	Gateway            string `json:"gateway"`                  // 网关IP地址
	CanUseHoneyIPRange string `json:"canUseHoneyIPRange"`       // 可用蜜罐IP范围
}

// UpdateView 处理网络信息更新的请求
func (NetApi) UpdateView(c *gin.Context) {
	// 获取并绑定更新请求的参数（包含网络ID、新标题等信息）
	cr := middleware.GetBind[UpdateRequest](c)

	// 查询指定ID的网络信息，用于后续更新和验证
	var model models.NetModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		// 若网络不存在，返回错误提示
		res.FailWithMsg("网络不存在", c)
		return
	}

	// 验证网络标题是否修改，若修改则检查新标题是否重复（排除当前网络ID）
	if cr.Title != model.Title {
		var newNet models.NetModel
		err = global.DB.Take(&newNet, "title = ? and id <> ?", cr.Title, cr.ID).Error
		if err == nil {
			// 若存在相同标题的其他网络，返回错误提示
			res.FailWithMsg("修改的网络名称不能重复", c)
			return
		}
	}

	// 验证网关IP（若网关不为空）
	if cr.Gateway != "" {
		// 解析网关IP是否有效
		gateway := net.ParseIP(cr.Gateway)
		if gateway == nil {
			res.FailWithMsg("网关ip格式错误", c)
			return
		}

		// 验证网关是否为IPv4
		ip4 := gateway.To4()
		if ip4 == nil {
			res.FailWithMsg("网关ip只支持ipv4", c)
			return
		}

		// 验证网关IP是否与探针IP相同（不允许相同）
		if cr.Gateway == model.IP {
			res.FailWithMsg("网关ip不能是探针ip", c)
			return
		}

		// 验证网关IP是否属于当前网络的子网
		_, _net, _ := net.ParseCIDR(model.Subnet())
		if !_net.Contains(gateway) {
			res.FailWithMsg("网关ip不属于当前子网", c)
			return
		}
	}

	// 验证可用蜜罐IP范围（若范围不为空）
	if cr.CanUseHoneyIPRange != "" {
		// 解析IP范围，获取具体IP列表
		ipList, err1 := ip.ParseIPRange(cr.CanUseHoneyIPRange)
		if err1 != nil {
			res.FailWithMsg(err1.Error(), c)
			return
		}

		// 检查每个IP是否属于当前网络的子网
		for _, s := range ipList {
			if !model.InSubnet(s) {
				res.FailWithMsg(fmt.Sprintf("%s不属于当前子网", s), c)
				return
			}
		}
	}

	// 执行网络信息更新操作
	err = global.DB.Model(&model).Updates(models.NetModel{
		Title:              cr.Title,
		Gateway:            cr.Gateway,
		CanUseHoneyIPRange: cr.CanUseHoneyIPRange,
	}).Error
	if err != nil {
		res.FailWithMsg("网络信息修改失败", c)
		return
	}

	// 更新成功，返回提示信息
	res.OkWithMsg("网络信息修改成功", c)
}
