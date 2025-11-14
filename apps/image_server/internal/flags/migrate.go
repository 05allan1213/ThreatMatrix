package flags

// File: flags/migrate.go
// Description: 实现 `Migrate` 命令，用于执行所有模型的表结构迁移。

import (
	"image_server/internal/global"
	"image_server/internal/models"

	"github.com/sirupsen/logrus"
)

// 迁移表结构
func Migrate() {
	// 检查数据库是否已初始化
	if global.DB == nil {
		logrus.Fatal("数据库未初始化，请检查程序执行流程")
	}

	err := global.DB.AutoMigrate(
		&models.HostTemplateModel{},   // 主机模板
		&models.ImageModel{},          // 镜像
		&models.MatrixTemplateModel{}, // 矩阵模板
		&models.ServiceModel{},        // 服务
	)
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
