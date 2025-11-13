package user_api

// File:api/user_api/login.go
// Description: 用户登录接口

import (
	"fmt"
	"honey_server/middleware"
	"honey_server/utils/res"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required" label:"用户名"`
	Password string `json:"password" binding:"required" label:"密码"`
}

// 用户登录接口
func (UserApi) LoginView(c *gin.Context) {
	cr := middleware.GetBind[LoginRequest](c)
	log := middleware.GetLog(c)
	log.Infof("这是请求的内容 %v", cr)
	log.Infof("ip %s", c.ClientIP())

	fmt.Println(cr)
	res.OkWithMsg("登录成功", c)
}
