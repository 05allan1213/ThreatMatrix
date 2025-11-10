package main

// File: main.go
// Description: 提供诱捕服务的应用入口逻辑。

import (
	"honey_server/core"
	"honey_server/flags"
	"honey_server/global"
)

func main() {
	global.Config = core.ReadConfig() // 读取配置文件
	global.Log = core.GetLogger()     // 初始化日志系统
	global.DB = core.InitDB()         // 初始化数据库
	flags.Run()                       // 解析命令行参数
	global.Log.Infof("info日志")
	global.Log.Warnf("warn日志")
	global.Log.Errorf("err日志")
}
