// Package global 暴露诱捕服务运行期间需要的全局变量。
//
// 本文件声明全局数据库实例指针，供其他模块共享。
package global

import "gorm.io/gorm"

var (
	DB *gorm.DB
)
