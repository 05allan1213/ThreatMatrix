package flags

// File: flags/migrate.go
// Description: 实现 `Migrate` 命令，用于执行所有模型的表结构迁移。

import (
	"honey_server/internal/global"
	"honey_server/internal/models"

	"github.com/sirupsen/logrus"
)

// 迁移表结构
func Migrate() {
	// 检查数据库是否已初始化
	if global.DB == nil {
		logrus.Fatal("数据库未初始化，请检查程序执行流程")
	}

	err := global.DB.AutoMigrate(
		&models.HoneyIpModel{},        // 诱捕IP
		&models.HoneyPortModel{},      // 诱捕端口
		&models.HostModel{},           // 存活主机
		&models.HostTemplateModel{},   // 主机模板
		&models.ImageModel{},          // 镜像
		&models.LogModel{},            // 日志
		&models.MatrixTemplateModel{}, // 矩阵模板
		&models.NetModel{},            // 网络
		&models.NodeModel{},           // 节点
		&models.NodeNetworkModel{},    // 节点网络
		&models.ServiceModel{},        // 服务
		&models.UserModel{},           // 用户
	)
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
