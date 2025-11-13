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
	Username string `json:"username"`
	Password string `json:"password"`
}

// 用户登录接口
func (UserApi) LoginView(c *gin.Context) {
	//var cr LoginRequest
	//err := c.ShouldBindJSON(&cr)
	//if err != nil {
	//	c.JSON(200, gin.H{"code": 7, "msg": "参数错误"})
	//	return
	//}
	cr := middleware.GetBind[LoginRequest](c)

	fmt.Println(cr)
	res.OkWithMsg("登录成功", c)
}
