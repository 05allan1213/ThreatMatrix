package user_api

// File: api/user_api/user_info.go
// Description: 用户信息相关接口，实现获取当前登录用户信息的功能

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// UserInfoResponse 用户信息返回结构体
type UserInfoResponse struct {
	UserID        uint   `json:"userID"`        // 用户ID
	Username      string `json:"username"`      // 用户名
	Role          int8   `json:"role"`          // 用户角色 1管理员 2普通用户
	LastLoginDate string `json:"lastLoginDate"` // 最近登录时间
}

// UserInfoView 获取当前登录用户信息
func (UserApi) UserInfoView(c *gin.Context) {
	auth := middleware.GetAuth(c) // 获取认证信息（包含用户ID）

	var user models.UserModel
	err := global.DB.Take(&user, auth.UserID).Error // 查询用户
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// 封装返回数据
	data := UserInfoResponse{
		UserID:        user.ID,
		Username:      user.Username,
		Role:          user.Role,
		LastLoginDate: user.LastLoginDate,
	}

	res.OkWithData(data, c)
}
