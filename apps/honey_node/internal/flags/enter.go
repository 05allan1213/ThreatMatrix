package flags

// File: flags/enter.go
// Description:
// 该模块实现一套可扩展的命令行指令体系，支持：
//   - 解析命令行参数
//   - 注册命令
//   - 输出帮助信息
//   - 执行预定义命令（如 db 迁移、打印版本、用户管理等）
// 整体结构支持“菜单 + 子命令”两级结构，并可通过 registerCommand 扩展。

import (
	"flag"
	"fmt"
	"honey_node/internal/global"
	"os"

	"github.com/sirupsen/logrus"
)

// FlagOptions 命令行参数结构体，用于存储所有解析后的命令行参数。
type FlagOptions struct {
	File    string // -f 配置文件路径
	Version bool   // -vv 是否输出版本信息
	DB      bool   // -db 是否执行数据库迁移
	Menu    string // -m 一级菜单
	Type    string // -t 二级命令
	Value   string // -v 子命令所需参数
	Help    bool   // -h 是否显示帮助
}

var Options FlagOptions // 全局命令行参数实例

// init 在程序启动阶段执行，负责初始化所有命令行参数并注册命令。
func init() {
	// 绑定命令行参数到 Options 结构体
	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件路径")
	flag.BoolVar(&Options.Version, "vv", false, "打印当前版本")
	flag.BoolVar(&Options.Help, "h", false, "帮助信息")
	flag.BoolVar(&Options.DB, "db", false, "迁移表结构")
	flag.StringVar(&Options.Menu, "m", "", "菜单 user")
	flag.StringVar(&Options.Type, "t", "", "类型 create list")
	flag.StringVar(&Options.Value, "v", "", "值")
	flag.Parse()

	// 注册所有可用命令
	RegisterCommand()
}

// RegisterCommand 注册所有一级菜单与子命令。
// 后续若有新的菜单或子命令，只需在这里添加 registerCommand 调用即可。
func RegisterCommand() {
}

// runBaseCommand 执行基础命令，其优先级最高。
// 若触发基础命令，则执行后直接退出程序。
func runBaseCommand() {
	// 执行数据库迁移
	if Options.DB {
		os.Exit(0)
	}

	// 输出版本信息
	if Options.Version {
		logrus.Infof("当前版本: %s  commit: %s, buildTime: %s",
			global.Version, global.Commit, global.BuildTime)
		os.Exit(0)
	}
}

// runHelpCommand 根据 -h 的输入输出帮助信息。
// 1. 没有输入菜单时，展示一级菜单列表。
// 2. 输入菜单但没有输入子菜单时，展示该菜单下所有子命令。
func runHelpCommand() {
	// ① -h（无菜单） → 展示所有一级菜单
	if Options.Menu == "" && Options.Type == "" && Options.Help {
		fmt.Printf("菜单项:\n")
		for key := range HelpCommandMap {
			fmt.Printf("%s 使用 -m %s -h 查看具体子菜单\n", key, key)
		}
		os.Exit(0)
	}

	// ② -m user -h → 展示 user 下所有子命令
	if Options.Menu != "" && Options.Type == "" && Options.Help {
		subMenuMap, ok := HelpCommandMap[Options.Menu]
		if !ok {
			logrus.Fatalf("不存在的菜单项 %s", Options.Menu)
		}
		for key, help := range subMenuMap {
			fmt.Printf("%s %s\n", key, help)
		}
		os.Exit(0)
	}
}

// runCommand 执行已经注册的命令。
// 逻辑：
//
//	输入 -m 与 -t → 根据 "m:t" 找到注册的 Command → 执行 Func()
func runCommand() {
	// 若二级命令不完整，则不执行
	if Options.Menu == "" || Options.Type == "" {
		return
	}

	key := fmt.Sprintf("%s:%s", Options.Menu, Options.Type)

	command, ok := CommandMap[key]
	if !ok {
		logrus.Fatalf("不存在的菜单项 %s %s", Options.Menu, Options.Type)
	}

	// 执行命令
	command.Func()
	os.Exit(0)
}

// 示例：
// ./main -db
// ./main -vv
// ./main -h
// ./main -m user -h
// ./main -m user -t list
// ./main -m user -t create
// ./main -m user -t create -v '{"username":"admin","password":"admin"}'

// Command 命令结构体。
// 一个完整命令由：菜单、子命令、帮助信息与执行函数组成。
type Command struct {
	Menu string // 一级菜单
	Type string // 子命令
	Help string // 帮助说明
	Func func() // 执行逻辑
}

var CommandMap = map[string]*Command{}              // “菜单:子命令” → Command
var HelpCommandMap = map[string]map[string]string{} // 菜单 → (子命令 → 帮助)

// registerCommand 注册单个命令。
// 同时维护：
//  1. CommandMap：执行映射表
//  2. HelpCommandMap：帮助映射表
func registerCommand(menu, subMenu, help string, fun func()) {
	key := fmt.Sprintf("%s:%s", menu, subMenu)

	// 存入执行映射（用于 runCommand）
	CommandMap[key] = &Command{
		Menu: menu,
		Type: subMenu,
		Help: help,
		Func: fun,
	}

	// 维护帮助映射
	subMenuMap, ok := HelpCommandMap[menu]
	if ok {
		subMenuMap[subMenu] = help
	} else {
		HelpCommandMap[menu] = map[string]string{
			subMenu: help,
		}
	}
}

// Run 程序入口命令执行流程。
// 执行顺序：
//  1. 基础命令（db、版本）
//  2. 帮助命令（-h）
//  3. 注册命令（用户命令等）
func Run() {
	runBaseCommand()
	runHelpCommand()
	runCommand()
}
