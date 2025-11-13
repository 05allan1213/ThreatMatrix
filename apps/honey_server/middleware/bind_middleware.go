package middleware

// File: middleware/bind_middleware.go
// Description: 通用请求参数绑定中间件

import (
	"honey_server/utils/res"

	"github.com/gin-gonic/gin"
)

// BindJsonMiddleware 通用 JSON 参数绑定中间件
func BindJsonMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	// 将绑定成功的结构体保存到上下文
	c.Set("request", cr)
}

// BindQueryMiddleware 通用 Query 参数绑定中间件
func BindQueryMiddleware[T any](c *gin.Context) {
	var cr T
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithMsg("参数绑定错误", c)
		c.Abort()
		return
	}
	c.Set("request", cr)
}

// GetBind 获取绑定参数的通用方法
func GetBind[T any](c *gin.Context) (cr T) {
	return c.MustGet("request").(T)
}
