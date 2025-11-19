package host_api

// File: api/host_api/list.go
// Description: 存活主机列表API

import (
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ListRequest 主机列表查询请求结构体
type ListRequest struct {
	models.PageInfo
	NodeID uint `form:"nodeID"` // 节点ID筛选条件
	NetID  uint `form:"netID"`  // 网络ID筛选条件
}

// ListResponse 主机列表查询响应结构体
type ListResponse struct {
	models.HostModel
	NetTitle  string `json:"netTitle"`  // 关联网络的名称
	NodeTitle string `json:"nodeTitle"` // 关联节点的名称
}

// ListView 主机列表查询接口处理函数
func (HostApi) ListView(c *gin.Context) {
	// 从请求中绑定并获取列表查询参数（包含分页和筛选条件）
	cr := middleware.GetBind[ListRequest](c)

	// 调用通用查询服务获取主机数据
	_list, count, _ := common_service.QueryList(models.HostModel{NodeID: cr.NodeID, NetID: cr.NetID}, common_service.QueryListRequest{
		Likes:    []string{"ip", "mac"},
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
		Preload:  []string{"NodeModel", "NetModel"},
	})

	// 将查询结果转换为包含关联名称的响应结构体列表
	var list = make([]ListResponse, 0)
	for _, model := range _list {
		list = append(list, ListResponse{
			HostModel: model,
			NodeTitle: model.NodeModel.Title, // 从关联的NodeModel获取节点名称
			NetTitle:  model.NetModel.Title,  // 从关联的NetModel获取网络名称
		})
	}

	// 返回分页列表数据（数据列表+总数）
	res.OkWithList(list, count, c)
}
