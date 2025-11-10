package global

// File: global/enter.go
// Description: 声明全局数据库实例指针，供其他模块共享。

import "gorm.io/gorm"

var (
	DB *gorm.DB
)
