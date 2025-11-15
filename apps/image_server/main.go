package main

// File: main.go
// Description: 提供诱捕服务的应用入口逻辑。

import (
	"image_server/internal/core"
	"image_server/internal/flags"
	"image_server/internal/global"
	"image_server/internal/routers"
	"image_server/internal/service/cron_service"
)

func main() {
	core.InitIPDB()                         // 初始化 IP 归属地数据库
	global.Config = core.ReadConfig()       // 读取配置文件
	core.SetLogDefault()                    // 设置日志默认配置
	global.DockerClient = core.InitDocker() // 初始化 Docker 客户端
	global.Log = core.GetLogger()           // 初始化日志系统
	global.DB = core.GetDB()                // 初始化数据库连接
	global.Redis = core.GetRedisClient()    // 初始化Redis连接
	flags.Run()                             // 解析命令行参数
	cron_service.Run()                      // 启动定时任务
	routers.Run()                           // 启动路由服务
}
