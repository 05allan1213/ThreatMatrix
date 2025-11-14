package main

// File: main.go
// Description: 提供诱捕服务的应用入口逻辑。

import (
	"honey_server/core"
	"honey_server/flags"
	"honey_server/global"
	"honey_server/routers"
)

func main() {
	core.InitIPDB()                      // 初始化 IP 归属地数据库
	global.Config = core.ReadConfig()    // 读取配置文件
	core.SetLogDefault()                 // 设置日志默认配置
	global.Log = core.GetLogger()        // 初始化日志系统
	global.DB = core.GetDB()             // 初始化数据库连接
	global.Redis = core.GetRedisClient() // 初始化Redis连接
	flags.Run()                          // 解析命令行参数
	routers.Run()                        // 启动路由服务
}
