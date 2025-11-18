package net_api

// File: api/net_api/list.go
// Description: 网络列表API

import (
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ListRequest 网络列表查询请求参数
type ListRequest struct {
	NodeID          uint `form:"nodeID"` // 节点ID
	models.PageInfo      // 分页信息
}

// ListResponse 网络列表查询响应结果
type ListResponse struct {
	models.NetModel        // 网络模型基础信息
	NodeTitle       string `json:"nodeTitle"`  // 所属节点的名称
	NodeStatus      int8   `json:"nodeStatus"` // 所属节点的状态
}

// ListView 处理网络列表查询的请求
func (NetApi) ListView(c *gin.Context) {
	// 获取并绑定列表查询的请求参数（包含节点ID和分页信息）
	cr := middleware.GetBind[ListRequest](c)

	// 调用通用查询服务查询网络列表
	_list, count, _ := common_service.QueryList(models.NetModel{NodeID: cr.NodeID}, common_service.QueryListRequest{
		Likes:    []string{"title", "ip"}, // 支持模糊查询的字段：名称、IP
		PageInfo: cr.PageInfo,             // 分页参数
		Sort:     "created_at desc",       // 排序方式：按创建时间倒序
		Preload:  []string{"NodeModel"},   // 预加载关联的节点模型信息
	})

	// 将查询到的网络列表转换为响应格式（补充节点名称和状态）
	var list = make([]ListResponse, 0)
	for _, model := range _list {
		_progress, ok := netProgressMap.Load(model.ID)
		if ok {
			progress := _progress.(float64)
			model.ScanProgress = progress
		}
		list = append(list, ListResponse{
			NetModel:   model,                  // 网络基础信息
			NodeTitle:  model.NodeModel.Title,  // 从关联的节点信息中获取节点名称
			NodeStatus: model.NodeModel.Status, // 从关联的节点信息中获取节点状态
		})
	}

	// 返回包含列表数据和总条数的成功响应
	res.OkWithList(list, count, c)
}
