package routers

// File: routers/honey_port_routers.go
// Description: 诱捕转发路由

import (
	"honey_server/internal/api"
	"honey_server/internal/api/honey_port_api"
	"honey_server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func HoneyPortRouters(r *gin.RouterGroup) {
	var app = api.App.HoneyPortApi

	// 诱捕转发更新（POST），绑定 JSON 请求体
	r.PUT("honey_port", middleware.BindJsonMiddleware[honey_port_api.UpdateRequest], app.UpdateView)

	// 诱捕转发列表（GET），绑定 Query 参数
	r.GET("honey_port", middleware.BindQueryMiddleware[honey_port_api.ListRequest], app.ListView)
}
