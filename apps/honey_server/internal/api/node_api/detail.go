package node_api

// File: api/node_api/detail.go
// Description: 节点详情接口

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// 根据ID获取节点详细信息
func (NodeApi) DetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NodeModel
	err := global.DB.Take(&model, cr.Id).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}
	res.OkWithData(model, c)
}
