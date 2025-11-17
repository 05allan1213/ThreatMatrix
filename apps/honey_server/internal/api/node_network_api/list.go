package node_network_api

// File: api/node_network_api/list.go
// Description: 提供节点网卡信息列表的API接口

import (
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ListRequest 节点网络列表请求参数结构体
type ListRequest struct {
	NodeID          uint `form:"nodeID" binding:"required"` // 节点ID
	models.PageInfo      // 嵌套分页信息结构体（包含Page页码、Limit每页数量）
}

// ListView 获取指定节点的网络信息列表
func (NodeNetworkApi) ListView(c *gin.Context) {
	// 获取并绑定请求参数（节点ID和分页信息），通过中间件进行参数校验
	cr := middleware.GetBind[ListRequest](c)

	// 调用通用查询服务，查询指定节点的网络信息列表
	list, count, _ := common_service.QueryList(
		models.NodeNetworkModel{NodeID: cr.NodeID}, // 查询的模型及固定条件（节点ID）
		common_service.QueryListRequest{
			Likes:    []string{"network", "ip"}, // 支持模糊搜索的字段（网卡名称、IP地址）
			PageInfo: cr.PageInfo,               // 分页参数（页码、每页数量）
			Sort:     "created_at desc",         // 排序方式：按创建时间降序
		},
	)

	// 返回成功响应，包含网络信息列表和总数
	res.OkWithList(list, count, c)
}
