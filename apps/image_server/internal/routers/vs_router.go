package routers

// File: routers/vs_router.go
// Description: 虚拟服务路由注册

import (
	"image_server/internal/api"
	"image_server/internal/api/vs_api"
	"image_server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func VsRouter(r *gin.RouterGroup) {
	app := api.App.VsApi

	// 虚拟服务创建（POST），绑定 JSON 请求体
	r.POST("vs", middleware.BindJsonMiddleware[vs_api.VsCreateRequest], app.VsCreateView)

	// 虚拟服务列表查询（GET），绑定 Query 参数
	r.GET("vs", middleware.BindQueryMiddleware[vs_api.VsListRequest], app.VsListView)

	// 虚拟服务选项列表查询（GET）
	r.GET("vs/options", app.VsOptionsListView)
}
