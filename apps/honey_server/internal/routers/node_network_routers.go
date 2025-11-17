package routers

// File: routers/node_network_routers.go
// Description: 节点网卡路由

import (
	"honey_server/internal/api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"

	"github.com/gin-gonic/gin"
)

// 节点网卡路由
func NodeNetworkRouters(r *gin.RouterGroup) {
	var app = api.App.NodeNetworkApi

	// 节点网卡刷新（GET），绑定 Query 参数
	r.GET("node_network/flush", middleware.BindQueryMiddleware[models.IDRequest], app.FlushView)
}
