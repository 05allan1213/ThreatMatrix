package node_api

// File: api/node_api/update.go
// Description: 提供节点信息更新的API接口

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// UpdateRequest 节点更新请求参数结构体
type UpdateRequest struct {
	ID    uint   `json:"id" binding:"required"`    // 节点ID
	Title string `json:"title" binding:"required"` // 节点新标题
}

// UpdateView 处理节点标题更新的API方法
func (NodeApi) UpdateView(c *gin.Context) {
	// 获取并绑定请求参数，通过中间件进行参数校验
	cr := middleware.GetBind[UpdateRequest](c)

	// 根据ID查询节点是否存在
	var model models.NodeModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}

	// 更新节点标题
	err = global.DB.Model(&model).Update("title", cr.Title).Error
	if err != nil {
		res.FailWithMsg("节点修改失败", c)
		return
	}

	// 返回更新成功的响应
	res.OkWithMsg("节点修改成功", c)
}
