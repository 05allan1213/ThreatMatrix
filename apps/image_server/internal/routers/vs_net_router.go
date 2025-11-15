package routers

// File: routers/vs_net_router.go
// Desciption: 虚拟网络路由

import (
	"image_server/internal/api"
	"image_server/internal/api/vs_net_api"
	"image_server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func VsNetRouter(r *gin.RouterGroup) {
	app := api.App.VsNetApi

	// 虚拟子网更新（PUT），绑定 JSON 请求体
	r.PUT("vs_net", middleware.BindJsonMiddleware[vs_net_api.VsNetRequest], app.VsNetUpdateView)

	// 虚拟子网列表（GET）
	r.GET("vs_net", app.VsNetInfoView)

}
