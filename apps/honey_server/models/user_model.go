package models

// File: models/user_model.go
// Description: 定义系统用户的数据模型，用于身份与权限管理。
import "gorm.io/gorm"

// 用户模型
type UserModel struct {
	gorm.Model
	Username      string `json:"username"`      // 用户名
	Role          string `json:"role"`          // 角色
	Password      string `json:"-"`             // 密码
	LastLoginDate string `json:"lastLoginDate"` // 最后登录时间
}
