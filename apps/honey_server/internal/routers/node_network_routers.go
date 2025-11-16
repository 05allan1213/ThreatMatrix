package routers

// File: routers/node_network_routers.go
// Description: 节点网卡路由

import (
	"honey_server/internal/api"

	"github.com/gin-gonic/gin"
)

// 节点网卡路由
func NodeNetworkRouters(r *gin.RouterGroup) {
	var app = api.App.NodeNetworkApi

	// 节点网卡刷新（GET）
	r.GET("node_network/flush", app.FlushView)
}
