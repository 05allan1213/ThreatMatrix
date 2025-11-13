package user_api

// File: api/user_api/create.go
// Description: 创建用户接口

import (
	"fmt"
	"honey_server/middleware"
	"honey_server/service/user_service"
	"honey_server/utils/res"

	"github.com/gin-gonic/gin"
)

// CreateRequest 定义创建用户时接收的请求参数结构体
type CreateRequest struct {
	Username string `json:"username" binding:"required" label:"用户名"` // 用户名（必填）
	Password string `json:"password" binding:"required" label:"密码"`  // 密码（必填）
	Role     int8   `json:"role" binding:"required,ne=1" label:"角色"` // 角色（必填，不能为1）
}

// CreateView 处理创建用户的 HTTP 请求
func (UserApi) CreateView(c *gin.Context) {
	// 从请求中绑定参数
	cr := middleware.GetBind[CreateRequest](c)

	// 获取请求上下文中的日志对象
	log := middleware.GetLog(c)

	// 创建用户服务实例
	us := user_service.NewUserService(log)

	// 调用用户服务的 Create 方法创建新用户
	user, err := us.Create(user_service.UserCreateRequest{
		Username: cr.Username,
		Password: cr.Password,
		Role:     cr.Role,
	})
	if err != nil {
		// 创建失败记录日志并返回错误响应
		msg := fmt.Sprintf("创建用户失败 %s", err)
		log.Errorf("%s", msg)
		res.FailWithMsg(msg, c)
		return
	}

	// 创建成功，返回用户ID
	res.OkWithData(user.ID, c)
}
