package core

// File: core/db.go
// Description: 定义数据库连接配置并实现 `InitDB`，用于建立与数据库的连接及配置连接池。

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 数据库配置
type DB struct {
	DbName   string `yaml:"db_name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// 初始化数据库
func InitDB() (database *gorm.DB) {
	var db = DB{
		DbName:   "honey_db",
		Host:     "82.157.155.26",
		Port:     3306,
		User:     "ayp",
		Password: "801026qwe",
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.DbName,
	)
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
	return
}
