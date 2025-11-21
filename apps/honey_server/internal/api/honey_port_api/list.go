package honey_port_api

// File: api/honey_port_api/list.go
// Description: 诱捕端口列表API

import (
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ListRequest 诱捕端口列表查询请求结构体
type ListRequest struct {
	models.PageInfo
	HoneyIPID uint `form:"honeyIpID" binding:"required"` // 关联的诱捕IP ID
}

// ListResponse 诱捕端口列表查询响应结构体
type ListResponse struct {
	models.HoneyPortModel
	ServiceTitle string `json:"serviceTitle"` // 关联服务的名称
}

// ListView 诱捕端口列表查询接口处理函数
func (HoneyPortApi) ListView(c *gin.Context) {
	// 从请求中绑定并获取列表查询参数（包含必填的诱捕IP ID）
	cr := middleware.GetBind[ListRequest](c)

	// 调用通用查询服务获取诱捕端口数据
	_list, count, _ := common_service.QueryList(models.HoneyPortModel{HoneyIpID: cr.HoneyIPID}, common_service.QueryListRequest{
		PageInfo: cr.PageInfo,
		Sort:     "created_at desc",
		Preload:  []string{"ServiceModel"},
	})

	// 将查询结果转换为包含关联服务名称的响应结构体列表
	var list = make([]ListResponse, 0)
	for _, model := range _list {
		list = append(list, ListResponse{
			HoneyPortModel: model,
			ServiceTitle:   model.ServiceModel.Title, // 从关联的ServiceModel获取服务名称
		})
	}

	// 返回分页列表数据（数据列表+总数）
	res.OkWithList(list, count, c)
}
