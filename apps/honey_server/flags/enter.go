// Package flags 负责解析并调度诱捕服务的命令行标志。
//
// 本文件定义全局命令行参数结构并在初始化阶段完成解析。
package flags

import (
	"flag"
	"os"
)

// 全局命令行参数结构
type FlagOptions struct {
	File    string // 配置文件路径
	Version bool   // 打印当前版本
	DB      bool   // 迁移表结构
}

var Options FlagOptions // 全局命令行参数实例

// 初始化命令行参数
func init() {
	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件路径")
	flag.BoolVar(&Options.Version, "v", false, "打印当前版本")
	flag.BoolVar(&Options.DB, "db", false, "迁移表结构")
	flag.Parse()
}

// 运行命令行参数
func Run() {
	if Options.DB {
		Migrate()
		os.Exit(0)
	}
}
