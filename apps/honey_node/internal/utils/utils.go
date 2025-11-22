package utils

// File: utils/utils.go
// Description: 提供通用的工具函数。

// InList 检查给定的键是否存在于列表中。
func InList[T comparable](list []T, key T) bool {
	for _, t := range list {
		if t == key {
			return true
		}
	}
	return false
}
