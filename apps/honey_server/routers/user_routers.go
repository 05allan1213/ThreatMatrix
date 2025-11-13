package routers

// File: routers/user_routers.go
// Description: 定义用户相关的路由。

import (
	"fmt"
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
}
