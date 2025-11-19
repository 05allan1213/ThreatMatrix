package routers

// File: routers/net_routers.go
// Description: 网络路由

import (
	"honey_server/internal/api"
	"honey_server/internal/api/net_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"

	"github.com/gin-gonic/gin"
)

func NetRouters(r *gin.RouterGroup) {
	var app = api.App.NetApi

	// 网络列表（GET），绑定 Query 参数
	r.GET("net", middleware.BindQueryMiddleware[net_api.ListRequest], app.ListView)

	// 网络选项（GET）
	r.GET("net/options", app.OptionsView)

	// 网络详情（GET），绑定 URI 参数
	r.GET("net/:id", middleware.BindUriMiddleware[models.IDRequest], app.DetailView)

	// 网络更新（PUT），绑定 JSON 参数
	r.PUT("net", middleware.BindJsonMiddleware[net_api.UpdateRequest], app.UpdateView)

	// 网络删除（DELETE），绑定 JSON 参数
	r.DELETE("net", middleware.BindJsonMiddleware[models.IDListRequest], app.RemoveView)

	// 网络扫描（POST），绑定 JSON 参数
	r.POST("net/scan", middleware.BindJsonMiddleware[models.IDRequest], app.ScanView)

	// 网络使用 IP 列表（GET），绑定 Query 参数
	r.GET("net/ip_list", middleware.BindQueryMiddleware[net_api.NetUseIPListRequest], app.NetUseIPListView)
}
