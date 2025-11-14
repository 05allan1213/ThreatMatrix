package routers

// File: routers/chptcha_routers.go
// Description: 图片验证码相关路由配置。

import (
	"honey_server/internal/api"

	"github.com/gin-gonic/gin"
)

// 图片验证码相关路由
func CaptchaRouters(r *gin.RouterGroup) {
	var app = api.App.CaptchaApi
	r.GET("captcha", app.GenerateView)
}
