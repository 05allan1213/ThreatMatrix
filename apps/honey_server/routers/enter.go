package routers

// File: routers/doc.go
// Description: 提供应用的路由配置和管理。

import (
	"honey_server/global"
	"honey_server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 启动路由服务
func Run() {
	system := global.Config.System
	gin.SetMode(system.Mode)
	r := gin.Default()                                         // 创建默认路由
	r.Static("uploads", "uploads")                             // 静态文件服务
	g := r.Group("honey_server")                               // 统一路由前缀 /honey_server
	g.Use(middleware.LogMiddleware, middleware.AuthMiddleware) // 系统必须登录才能访问

	UserRouters(g)    // 用户相关路由
	CaptchaRouters(g) // 图片验证码路由

	webAddr := system.WebAddr
	logrus.Infof("web addr run %s", webAddr)

	r.Run(webAddr)
}
