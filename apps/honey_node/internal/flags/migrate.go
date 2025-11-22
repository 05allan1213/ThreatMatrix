package flags

// File: flags/migrate.go
// Description: 数据库表结构迁移

import (
	"honey_node/internal/global"
	"honey_node/internal/models"

	"github.com/sirupsen/logrus"
)

// 数据库表结构迁移
func Migrate() {
	err := global.DB.AutoMigrate(
		&models.PortModel{},
		&models.IpModel{},
	)
	if err != nil {
		logrus.Fatalf("表结构迁移失败 %s", err)
	}
	logrus.Infof("表结构迁移成功")
}
