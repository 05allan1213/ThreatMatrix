package routers

// File: routers/user_routers.go
// Description: 定义用户相关的路由。

import (
	"fmt"
	"honey_server/api"
	"honey_server/api/user_api"
	"honey_server/middleware"

	"github.com/gin-gonic/gin"
)

// 定义用户相关路由
func UserRouters(r *gin.RouterGroup) {
	r.GET("users", func(c *gin.Context) {
		fmt.Println(middleware.GetAuth(c))
		c.JSON(200, gin.H{"code": 0})
	})

	r.GET("login", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 1})
	})

	var app = api.App.UserApi
	r.POST("login", middleware.BindJsonMiddleware[user_api.LoginRequest], app.LoginView)
	r.POST("users", middleware.AdminMiddleware, middleware.BindJsonMiddleware[user_api.CreateRequest], app.CreateView)
}
