package node_api

// File: api/node_api/list.go
// Description: 节点列表api

import (
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// 获取节点列表
func (NodeApi) ListView(c *gin.Context) {
	cr := middleware.GetBind[models.PageInfo](c)
	list, count, _ := common_service.QueryList(models.NodeModel{}, common_service.QueryListRequest{
		Likes:    []string{"title", "ip"},
		PageInfo: cr,
		Sort:     "created_at desc",
	})
	res.OkWithList(list, count, c)
}
