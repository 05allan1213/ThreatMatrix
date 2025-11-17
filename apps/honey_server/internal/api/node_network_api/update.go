package node_network_api

// File: api/node_network_api/update.go
// Description: 提供节点网卡网关信息更新的API接口

import (
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"
	"net"

	"github.com/gin-gonic/gin"
)

// UpdateRequest 节点网卡更新请求参数结构体
type UpdateRequest struct {
	ID      uint   `json:"id" binding:"required"` // 网卡ID
	Gateway string `json:"gateway"`               // 网关地址
}

// UpdateView 处理节点网卡网关更新的API方法
func (NodeNetworkApi) UpdateView(c *gin.Context) {
	// 获取并绑定请求参数，通过中间件进行参数校验
	cr := middleware.GetBind[UpdateRequest](c)

	// 根据ID查询网卡记录是否存在
	var model models.NodeNetworkModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("节点网卡不存在", c)
		return
	}

	// 当网关地址不为空时，执行网关验证逻辑
	if cr.Gateway != "" {
		// 1. 验证网关IP格式是否合法
		gateway := net.ParseIP(cr.Gateway)
		if gateway == nil {
			res.FailWithMsg("网关ip格式错误", c)
			return
		}

		// 2. 验证网关是否为IPv4地址（不支持IPv6）
		ip4 := gateway.To4()
		if ip4 == nil {
			res.FailWithMsg("网关ip只支持ipv4", c)
			return
		}

		// 3. 验证网关IP不能与当前网卡的IP（探针IP）相同
		if cr.Gateway == model.IP {
			res.FailWithMsg("网关ip不能是探针ip", c)
			return
		}

		// 4. 验证网关IP是否属于当前网卡所在的子网
		// 构建当前网卡的子网CIDR（格式：IP/掩码长度）
		_, _net, _ := net.ParseCIDR(fmt.Sprintf("%s/%d", model.IP, model.Mask))
		// 检查网关是否在当前子网内
		if !_net.Contains(gateway) {
			res.FailWithMsg("网关ip不属于当前子网", c)
			return
		}
	}

	// 执行数据库更新操作，更新网关字段
	err = global.DB.Model(&model).Update("gateway", cr.Gateway).Error
	if err != nil {
		res.FailWithMsg("节点网卡修改失败", c)
		return
	}

	// 返回更新成功的响应
	res.OkWithMsg("节点网卡修改成功", c)
}
