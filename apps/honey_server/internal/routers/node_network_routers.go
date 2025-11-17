package routers

// File: routers/node_network_routers.go
// Description: 节点网卡路由

import (
	"honey_server/internal/api"
	"honey_server/internal/api/node_network_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"

	"github.com/gin-gonic/gin"
)

// 节点网卡路由
func NodeNetworkRouters(r *gin.RouterGroup) {
	var app = api.App.NodeNetworkApi

	// 节点网卡刷新（GET），绑定 Query 参数
	r.GET("node_network/flush", middleware.BindQueryMiddleware[models.IDRequest], app.FlushView)

	// 节点网卡列表（GET），绑定 Query 参数
	r.GET("node_network", middleware.BindQueryMiddleware[node_network_api.ListRequest], app.ListView)

	// 节点网卡更新（PUT），绑定 JSON 参数
	r.PUT("node_network", middleware.BindJsonMiddleware[node_network_api.UpdateRequest], app.UpdateView)

	// 节点网卡启用（PUT），绑定 JSON 参数
	r.PUT("node_network/enable", middleware.BindJsonMiddleware[models.IDRequest], app.EnableView)

	// 节点网卡删除（DELETE），绑定 Uri 参数
	r.DELETE("node_network/:id", middleware.BindUriMiddleware[models.IDRequest], app.RemoveView)

}
