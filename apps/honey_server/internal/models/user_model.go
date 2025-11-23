package models

import "gorm.io/gorm"

// File: models/user_model.go
// Description: 定义系统用户的数据模型，用于身份与权限管理。

// 用户模型
type UserModel struct {
	Model
	Username      string `gorm:"size:32;index:idx_username" json:"username"` // 用户名
	Role          int8   `json:"role"`                                       // 角色 1:管理员 2:普通用户
	Password      string `gorm:"size:64" json:"-"`                           // 密码
	LastLoginDate string `gorm:"size:32" json:"lastLoginDate"`               // 最后登录时间
}

func (UserModel) BeforeDelete(tx *gorm.DB) error {
	return nil
}
