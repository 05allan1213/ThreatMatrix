package vs_api

// File: api/vs_api/vs_list.go
// Description: 虚拟服务列表 API，负责查询虚拟服务列表，支持分页、按端口/IP/标题筛选。

import (
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// VsListRequest 查询虚拟服务列表的请求参数结构
type VsListRequest struct {
	models.PageInfo        // 嵌入分页信息（包含页码、每页条数等）
	Port            int    `form:"port"`  // 筛选条件：服务端口
	IP              string `form:"ip"`    // 筛选条件：服务IP地址
	Title           string `form:"title"` // 筛选条件：服务标题（支持模糊查询）
}

// VsListView 获取虚拟服务列表的API入口函数
func (VsApi) VsListView(c *gin.Context) {
	// 绑定并验证请求参数（从上下文获取VsListRequest结构体，包含分页和筛选条件）
	cr := middleware.GetBind[VsListRequest](c)

	// 调用通用查询服务查询虚拟服务列表
	// 1. 构建查询条件：根据请求参数中的标题、IP、端口筛选
	// 2. 配置查询参数：标题支持模糊查询、指定分页信息、按创建时间倒序排序
	list, count, _ := common_service.QueryList(models.ServiceModel{
		Title: cr.Title,
		IP:    cr.IP,
		Port:  cr.Port,
	},
		common_service.QueryListRequest{
			Likes:    []string{"title"}, // 标题字段支持模糊查询
			PageInfo: cr.PageInfo,       // 分页参数
			Sort:     "created_at desc", // 排序规则：按创建时间倒序
		})

	// 返回查询结果（包含列表数据和总条数）
	res.OkWithList(list, count, c)
}
