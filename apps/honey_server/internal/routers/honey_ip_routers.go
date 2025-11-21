package routers

// File: routers/honey_ip_routers.go
// Description: 诱捕IP路由

import (
	"honey_server/internal/api"
	"honey_server/internal/api/honey_ip_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"

	"github.com/gin-gonic/gin"
)

func HoneyIPRouters(r *gin.RouterGroup) {
	var app = api.App.HoneyIPApi

	// 诱捕IP创建（POST），绑定 JSON 请求体
	r.POST("honey_ip", middleware.BindJsonMiddleware[honey_ip_api.CreateRequest], app.CreateView)

	// 诱捕IP列表查询（GET），绑定 Query 参数
	r.GET("honey_ip", middleware.BindQueryMiddleware[honey_ip_api.ListRequest], app.ListView)

	// 诱捕IP删除（DELETE），绑定 JSON 请求体
	r.DELETE("honey_ip", middleware.BindJsonMiddleware[models.IDListRequest], app.RemoveView)
}
