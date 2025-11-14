package routers

// File: routers/log_routers.go
// Description: 日志模块路由注册，负责绑定日志查询与删除接口

import (
	"honey_server/internal/api"
	"honey_server/internal/api/log_api"
	"honey_server/internal/middleware"
	"honey_server/internal/models"

	"github.com/gin-gonic/gin"
)

// LogRouters 注册日志模块相关路由
func LogRouters(r *gin.RouterGroup) {
	var app = api.App.LogApi // 获取日志 API 实例

	// 日志列表：GET + Query 绑定
	r.GET("logs",
		middleware.AdminMiddleware,                             // 管理员验证
		middleware.BindQueryMiddleware[log_api.LogListRequest], // 绑定查询参数
		app.LogListView,                                        // 处理请求
	)

	// 删除日志：DELETE + JSON 参数绑定
	r.DELETE("logs",
		middleware.AdminMiddleware,                          // 管理员验证
		middleware.BindJsonMiddleware[models.IDListRequest], // 绑定 JSON ID 列表
		app.RemoveView, // 处理请求
	)
}
