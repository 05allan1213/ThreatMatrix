package routers

// File: routers/doc.go
// Description: 提供应用的路由配置和管理。

import (
	"image_server/internal/global"
	"image_server/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 启动路由服务
func Run() {
	system := global.Config.System
	gin.SetMode(system.Mode)
	r := gin.Default() // 创建默认路由
	//r.Static("uploads", "uploads")                             // 静态文件服务
	g := r.Group("image_server")                               // 统一路由前缀 /honey_server
	g.Use(middleware.LogMiddleware, middleware.AuthMiddleware) // 系统必须登录才能访问，所有以 /honey_server 开头的路由默认都需要认证

	MirrorCloudRouter(g) // 镜像云相关路由
	VsRouter(g)          // 虚拟服务相关路由

	webAddr := system.WebAddr
	logrus.Infof("web addr run %s", webAddr)

	r.Run(webAddr)
}
