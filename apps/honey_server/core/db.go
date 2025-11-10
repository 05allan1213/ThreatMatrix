package core

// File: core/db.go
// Description: 实现数据库初始化，用于建立与数据库的连接及配置连接池。

import (
	"honey_server/global"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库
func InitDB() (database *gorm.DB) {
	dsn := global.Config.DB.DSN()
	dialector := mysql.Open(dsn)
	database, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 不生成实体外键
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s", err)
		return
	}
	// 配置连接池
	sqlDB, err := database.DB()
	if err != nil {
		logrus.Fatalf("获取数据库连接失败 %s", err)
		return
	}
	err = sqlDB.Ping()
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s", err)
		return
	}
	// 设置连接池
	sqlDB.SetMaxIdleConns(10)           // 设置连接池最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 设置连接池最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 设置连接池连接最大生命周期

	logrus.Infof("数据库连接成功")
	return
}
