package routers

// File: routers/host_routers.go
// Description: 存活主机路由

import (
	"honey_server/internal/api"
	"honey_server/internal/api/host_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"

	"github.com/gin-gonic/gin"
)

func HostRouters(r *gin.RouterGroup) {
	var app = api.App.HostApi

	// 存活主机列表查询（GET），绑定 Query 参数
	r.GET("host", middleware.BindQueryMiddleware[host_api.ListRequest], app.ListView)

	// 删除主机（DELETE），绑定 JSON 参数
	r.DELETE("host", middleware.BindJsonMiddleware[models.IDListRequest], app.RemoveView)
}
