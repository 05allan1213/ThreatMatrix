package routers

// File: routers/host_template_router.go
// Description: 主机模板路由

import (
	"image_server/internal/api"
	"image_server/internal/api/host_template_api"
	"image_server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func HostTemplateRouter(r *gin.RouterGroup) {
	app := api.App.HostTemplateApi

	// 主机模板创建（POST），绑定 JSON 请求体
	r.POST("host_template", middleware.BindJsonMiddleware[host_template_api.CreateRequest], app.CreateView)
}
