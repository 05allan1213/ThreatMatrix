package routers

// File: routers/doc.go
// Description: 提供应用的路由配置和管理。

import (
	"honey_server/global"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 启动路由服务
func Run() {
	system := global.Config.System
	gin.SetMode(system.Mode)
	r := gin.Default()
	g := r.Group("honey_server") // 统一路由前缀 /honey_server
	g.Use()                      // 添加中间件

	UserRouters(g)

	webAddr := system.WebAddr
	logrus.Infof("web addr run %s", webAddr)

	r.Run(webAddr)
}
