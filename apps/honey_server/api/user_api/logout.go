package user_api

// File: user_logout.go
// Description: 用户注销接口，记录日志并返回注销成功信息。

import (
	"honey_server/middleware"
	"honey_server/utils/res"
	"time"

	"github.com/gin-gonic/gin"
)

// UserLogoutView 用户注销接口
func (UserApi) UserLogoutView(c *gin.Context) {
	token := c.GetHeader("token")             // 获取用户请求头内的 token
	log := middleware.GetLog(c)               // 获取请求日志记录器
	auth := middleware.GetAuth(c)             // 获取当前用户认证信息
	expiresAt := time.Unix(auth.ExpiresAt, 0) // token 过期时间

	// 输出注销日志：用户ID、token、过期时间
	log.Infof("用户注销 %d %s %s", auth.UserID, token, expiresAt)

	// 返回成功响应
	res.OkWithMsg("注销成功", c)
}
