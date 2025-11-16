package routers

// File: routers/node_routers.go
// Description: 节点路由注册

import (
	"honey_server/internal/api"
	"honey_server/internal/api/node_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"

	"github.com/gin-gonic/gin"
)

func NodeRouters(r *gin.RouterGroup) {
	var app = api.App.NodeApi

	// 节点列表（GET），绑定 Query 参数
	r.GET("node", middleware.BindQueryMiddleware[models.PageInfo], app.ListView)

	// 节点详情（GET），绑定 URI 参数
	r.GET("node/:id", middleware.BindUriMiddleware[models.IDRequest], app.DetailView)

	// 节点更新（PUT），绑定 JSON 请求体
	r.PUT("node", middleware.BindJsonMiddleware[node_api.UpdateRequest], app.UpdateView)

	// 节点选项（GET）
	r.GET("node/options", app.OptionsView)

	// 节点删除（DELETE），绑定 URI 参数
	r.DELETE("node/:id", middleware.BindUriMiddleware[models.IDRequest], app.RemoveView)
}
