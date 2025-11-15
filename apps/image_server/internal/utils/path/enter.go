package path

// File: utils/path/enter.go
// Description: 项目路径工具函数，提供获取当前程序运行根路径的能力

import (
	"os"
)

// GetRootPath 获取项目当前运行的根目录路径
// 若获取失败则返回空字符串
func GetRootPath() (path string) {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}
