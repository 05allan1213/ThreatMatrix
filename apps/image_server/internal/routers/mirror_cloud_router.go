package routers

// File: routers/mirror_cloud_router.go
// Description: 镜像云模块路由注册

import (
	"image_server/internal/api"
	"image_server/internal/api/mirror_cloud_api"
	"image_server/internal/middleware"

	"github.com/gin-gonic/gin"
)

// MirrorCloudRouter 注册镜像云相关路由
func MirrorCloudRouter(r *gin.RouterGroup) {
	app := api.App.MirrorCloudApi // 获取镜像云 API 实例

	// 镜像文件查看（POST）
	r.POST("mirror_cloud/see", app.ImageSeeView)

	// 镜像创建（POST），绑定 JSON 请求体
	r.POST("mirror_cloud", middleware.BindJsonMiddleware[mirror_cloud_api.ImageCreateRequest], app.ImageCreateView)

	// 镜像列表查询（GET），绑定 Query 参数
	r.GET("mirror_cloud", middleware.BindQueryMiddleware[mirror_cloud_api.ImageListRequest], app.ImageListView)
}
