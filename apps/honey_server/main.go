package main

// File: main.go
// Description: 提供诱捕服务的应用入口逻辑。

import (
	"honey_server/core"
	"honey_server/flags"
	"honey_server/global"

	"github.com/sirupsen/logrus"
)

func main() {
	global.Config = core.ReadConfig() // 读取配置文件
	if global.Config == nil {
		logrus.Fatalf("配置文件读取失败")
		return
	}
	global.DB = core.InitDB() // 初始化数据库
	if global.DB == nil {
		logrus.Fatalf("数据库初始化失败")
		return
	}
	flags.Run() // 解析命令行参数
}
