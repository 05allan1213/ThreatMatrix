package honey_ip_api

// File: api/honey_ip_api/remove.go
// Description: 诱捕IP删除API

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// RemoveView 诱捕IP批量删除接口处理函数
func (HoneyIPApi) RemoveView(c *gin.Context) {
	// 从请求中绑定并获取批量删除的ID列表参数（models.IDListRequest结构体）
	cr := middleware.GetBind[models.IDListRequest](c)

	// 根据ID列表查询对应的诱捕IP记录，并预加载关联的节点模型（NodeModel）
	var honeyIPList []models.HoneyIpModel
	global.DB.Preload("NodeModel").Find(&honeyIPList, "id in ?", cr.IdList)

	// 未查询到任何诱捕IP记录则返回错误
	if len(honeyIPList) == 0 {
		res.FailWithMsg("未找到诱捕ip", c)
		return
	}

	// 获取第一条记录关联的节点模型（假设批量删除的诱捕IP属于同一节点）
	nodeModel := honeyIPList[0].NodeModel

	// 合法性校验1：判断节点状态是否为运行中
	if nodeModel.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 合法性校验2：通过gRPC检查节点是否在线
	_, ok := grpc_service.GetNodeCommand(nodeModel.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 下发批量删除的任务到消息队列（异步处理实际删除逻辑）

	// 更新诱捕IP状态为删除中（状态码4）
	global.DB.Model(&honeyIPList).Update("status", 4)

	// 返回删除任务启动成功的提示
	res.OkWithMsg("批量删除中", c)
}
