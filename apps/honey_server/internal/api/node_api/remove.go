package node_api

// File: api/node_api/remove.go
// Description: 提供节点删除的API接口

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// 处理节点删除
func (NodeApi) RemoveView(c *gin.Context) {
	// 获取并绑定请求参数（节点ID），通过中间件进行参数处理
	cr := middleware.GetBind[models.IDRequest](c)

	// 根据ID查询节点是否存在
	var model models.NodeModel
	err := global.DB.Take(&model, cr.Id).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}

	// 执行节点删除操作
	err = global.DB.Delete(&model).Error
	if err != nil {
		res.FailWithMsg("节点删除失败", c)
		return
	}

	// 返回删除成功的响应
	res.OkWithMsg("节点删除成功", c)
}
