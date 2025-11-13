package user_service

// File: service/user_service/enter.go
// Description: 用户服务，提供用户相关的业务逻辑处理

import "github.com/sirupsen/logrus"

// UserService 用户服务结构体
type UserService struct {
	log *logrus.Entry
}

// 创建用户服务实例
func NewUserService(log *logrus.Entry) *UserService {
	return &UserService{
		log: log,
	}
}
