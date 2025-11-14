package routers

// File: routers/mirror_cloud_router.go
// Description: 镜像云模块路由注册，负责绑定镜像查看相关接口

import (
	"image_server/internal/api"

	"github.com/gin-gonic/gin"
)

// MirrorCloudRouter 注册镜像云相关路由
func MirrorCloudRouter(r *gin.RouterGroup) {
	app := api.App.MirrorCloudApi // 获取镜像云 API 实例

	// 镜像文件查看接口
	r.POST("mirror_cloud/see", app.ImageSeeView)
}
