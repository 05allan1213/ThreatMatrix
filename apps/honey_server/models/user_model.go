// Package models 定义诱捕服务所使用的数据实体。
//
// 本文件描述系统用户的信息结构，用于身份与权限管理。
package models

import "gorm.io/gorm"

// 用户模型
type UserModel struct {
	gorm.Model
	Username      string `json:"username"`      // 用户名
	Role          string `json:"role"`          // 角色
	Password      string `json:"-"`             // 密码
	LastLoginDate string `json:"lastLoginDate"` // 最后登录时间
}
