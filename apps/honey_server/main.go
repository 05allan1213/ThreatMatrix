// Package main 提供诱捕服务的应用入口逻辑。
//
// 本文件负责初始化全局数据库连接并触发命令行子命令的执行。
package main

import (
	"honey_server/core"
	"honey_server/flags"
	"honey_server/global"
)

func main() {
	global.DB = core.InitDB()
	flags.Run()
}
