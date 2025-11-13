package user_service

// File: service/user_service/user_create.go
// Description: 实现用户创建的业务逻辑

import (
	"fmt"
	"honey_server/global"
	"honey_server/models"
	"honey_server/utils/pwd"
)

// UserCreateRequest 用户创建请求结构体
type UserCreateRequest struct {
	Role     int8   `json:"role"`     // 用户角色（1管理员 2普通用户）
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// Create 创建新用户
func (u *UserService) Create(req UserCreateRequest) (user models.UserModel, err error) {
	// 检查用户名是否已存在
	err = global.DB.Take(&user, "username = ?", req.Username).Error
	if err == nil {
		err = fmt.Errorf("%s 用户名已存在", req.Username)
		return
	}

	// 对用户密码进行加密
	hashPwd, _ := pwd.GenerateFromPassword(req.Password)

	// 构造新用户模型
	user = models.UserModel{
		Username: req.Username,
		Password: hashPwd,
		Role:     req.Role,
	}

	// 将用户信息写入数据库
	err = global.DB.Create(&user).Error
	if err != nil {
		err = fmt.Errorf("用户创建失败 %s", err)
		return
	}

	// 创建成功记录日志
	u.log.Infof("%s 用户创建成功", req.Username)
	return
}
