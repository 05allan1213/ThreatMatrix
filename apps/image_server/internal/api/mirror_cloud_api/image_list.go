package mirror_cloud_api

// File: api/mirror_cloud_api/image_list.go
// Description: 镜像列表接口，支持分页、模糊查询与通用数据查询服务

import (
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ImageListRequest 镜像列表请求结构体，包含分页参数
type ImageListRequest struct {
	models.PageInfo
}

// ImageListView 镜像列表查询接口
// 支持 title、image_name 的模糊搜索；
// 使用通用 QueryList 实现分页、排序与条件查询
func (MirrorCloudApi) ImageListView(c *gin.Context) {
	cr := middleware.GetBind[ImageListRequest](c)

	list, count, _ := common_service.QueryList(
		models.ImageModel{},
		common_service.QueryListRequest{
			Likes:    []string{"title", "image_name"}, // 模糊查询字段
			PageInfo: cr.PageInfo,                     // 分页信息
			Sort:     "created_at desc",               // 按创建时间倒序
		},
	)

	res.OkWithList(list, count, c)
}
