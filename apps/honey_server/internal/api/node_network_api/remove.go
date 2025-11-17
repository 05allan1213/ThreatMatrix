package node_network_api

// File: api/node_network_api/remover.go
// Description: 节点网络相关API接口处理逻辑，包含网卡删除等操作的实现

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// RemoveView 处理删除网卡的请求
func (NodeNetworkApi) RemoveView(c *gin.Context) {
	// 从请求中获取绑定的ID参数（用于指定要删除的网卡）
	cr := middleware.GetBind[models.IDRequest](c)

	// 查询指定ID的网卡信息
	var model models.NodeNetworkModel
	err := global.DB.Take(&model, cr.Id).Error
	if err != nil {
		// 若查询失败（网卡不存在），返回错误提示
		res.FailWithMsg("网卡不存在", c)
		return
	}

	// 执行网卡信息删除操作
	err = global.DB.Delete(&model).Error
	if err != nil {
		// 若删除失败，返回错误提示并附带具体错误信息
		res.FailWithMsg("网卡删除失败"+err.Error(), c)
		return
	}

	// 删除成功，返回成功提示
	res.OkWithMsg("网卡删除成功", c)
	return
}
