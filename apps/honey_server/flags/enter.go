package flags

// File: flags/enter.go
// Description: 定义全局命令行参数结构并在初始化阶段完成解析。

import (
	"flag"
	"honey_server/global"
	"os"

	"github.com/sirupsen/logrus"
)

// 全局命令行参数结构
type FlagOptions struct {
	File    string // 配置文件路径
	Version bool   // 打印当前版本
	DB      bool   // 迁移表结构
	Menu    string // 初始化菜单
	Type    string // 初始化类型
	Value   string // 初始化值
}

var Options FlagOptions // 全局命令行参数实例

// 初始化命令行参数
func init() {
	flag.StringVar(&Options.File, "f", "settings.yaml", "配置文件路径")
	flag.BoolVar(&Options.Version, "vv", false, "打印当前版本")
	flag.BoolVar(&Options.DB, "db", false, "迁移表结构")
	flag.StringVar(&Options.Menu, "m", "", "菜单 user")
	flag.StringVar(&Options.Type, "t", "", "类型 create list")
	flag.StringVar(&Options.Value, "v", "", "值")
	flag.Parse()
}

// 运行命令行参数
func Run() {
	// 若命令行参数中包含 -db，则执行数据库迁移操作并退出程序
	if Options.DB {
		Migrate()
		os.Exit(0)
	}

	// 若命令行参数中包含 -vv，则输出当前版本信息（版本号、提交号、构建时间）后退出程序
	if Options.Version {
		logrus.Infof("当前版本: %s  commit: %s, buildTime: %s",
			global.Version, global.Commit, global.BuildTime)
		os.Exit(0)
	}

	// 根据命令行参数的主菜单项进行分支处理
	switch Options.Menu {
	case "user": // 用户相关操作菜单
		var user User
		switch Options.Type {
		case "create": // 创建用户
			user.Create()
			os.Exit(0)
		case "list": // 查看用户列表
			user.List()
			os.Exit(0)
		default:
			// 若子菜单项不正确，打印错误日志并退出
			logrus.Fatalf("用户子菜单项不正确")
		}

	case "": // 若未指定菜单项则继续执行主程序（不退出）
		// 空分支，表示不执行额外命令行逻辑

	default:
		// 输入了不存在的菜单项，打印错误并终止程序
		logrus.Fatalf("菜单项不正确")
	}
}
