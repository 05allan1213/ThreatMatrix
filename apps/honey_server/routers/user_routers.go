package routers

// File: routers/user_routers.go
// Description: 定义用户相关的路由。

import (
	"honey_server/api"
	"honey_server/api/user_api"
	"honey_server/middleware"

	"github.com/gin-gonic/gin"
)

// 定义用户相关路由
func UserRouters(r *gin.RouterGroup) {
	var app = api.App.UserApi

	// 用户登录（POST），绑定 JSON 请求体
	r.POST("login", middleware.BindJsonMiddleware[user_api.LoginRequest], app.LoginView)

	// 创建用户（POST），管理员权限 + JSON 请求体绑定
	r.POST("users", middleware.AdminMiddleware, middleware.BindJsonMiddleware[user_api.CreateRequest], app.CreateView)

	// 用户列表查询（GET），绑定 Query 参数
	r.GET("users", middleware.BindQueryMiddleware[user_api.UserListRequest], app.UserListView)

	// 用户注销（POST）
	r.POST("logout", app.UserLogoutView)

	// 删除用户（DELETE），绑定 JSON 请求体
	r.DELETE("users", middleware.BindJsonMiddleware[user_api.UserRemoveRequest], app.UserRemoveView)

	// 获取用户信息（GET）
	r.GET("users/info", app.UserInfoView)
}
