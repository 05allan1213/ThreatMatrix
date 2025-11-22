package core

// File: core/db.go
// Description: 负责SQLite数据库连接的初始化、连接池配置及单例模式管理，提供全局唯一的数据库实例

import (
	"honey_node/internal/global"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// InitDB 初始化SQLite数据库连接并配置连接池
func InitDB() (db *gorm.DB) {
	cfg := global.Config.DB // 获取全局配置中的数据库配置项

	// 打开SQLite数据库连接（文件为gorm.db），禁用外键约束迁移以提升兼容性
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 不自动生成实体外键约束
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s", err)
		return
	}

	// 获取底层sql.DB实例，用于配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("获取数据库连接失败 %s", err)
		return
	}

	// 测试数据库连接是否可用
	err = sqlDB.Ping()
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s", err)
		return
	}

	// 设置连接池默认参数（若配置未指定则使用默认值）
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 10 // 默认最大空闲连接数
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 100 // 默认最大打开连接数
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = 10000 // 默认连接最大存活时间（秒）
	}

	// 打印连接池配置信息，便于调试
	logrus.Infof("最大空闲数 %d", cfg.MaxIdleConns)
	logrus.Infof("最大连接数 %d", cfg.MaxOpenConns)
	logrus.Infof("超时时间 %s", time.Duration(cfg.ConnMaxLifetime)*time.Second)

	// 配置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)                                    // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)                                    // 设置最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second) // 设置连接最大存活时间

	logrus.Infof("数据库连接成功")
	return db
}

// 全局数据库实例及单例锁：确保应用生命周期内只有一个数据库连接实例
var db *gorm.DB
var onceMysql sync.Once // 用于保证InitDB只执行一次

// GetDB 获取全局唯一的数据库实例（单例模式）
func GetDB() *gorm.DB {
	onceMysql.Do(func() {
		db = InitDB() // 首次调用时初始化数据库连接
	})
	return db
}
