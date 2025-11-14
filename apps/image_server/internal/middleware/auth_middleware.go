package middleware

// File: middleware/auth_middleware.go
// Description: 提供认证相关的中间件功能。

import (
	"image_server/internal/global"
	"image_server/internal/utils"
	"image_server/internal/utils/jwts"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 通用认证中间件
func AuthMiddleware(c *gin.Context) {
	// 先判断这个路径在不在白名单中
	path := c.Request.URL.Path
	if utils.InList(global.Config.WhiteList, path) {
		// 在白名单中，直接放行（免认证）
		c.Next()
		return
	}
	token := c.GetHeader("token") // 从请求头获取 token
	claims, err := jwts.ParseToken(token)
	if err != nil {
		// token 无效或解析失败
		res.FailWithMsg("认证失败", c)
		c.Abort() // 阻止后续处理函数执行
		return
	}

	// 解析成功，继续执行下一个中间件或处理函数
	c.Set("claims", claims) // 将解析后的用户信息存入上下文，供后续使用
	c.Next()
}

// GetAuth 获取当前请求的用户信息
func GetAuth(c *gin.Context) *jwts.Claims {
	return c.MustGet("claims").(*jwts.Claims)
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware(c *gin.Context) {
	token := c.GetHeader("token") // 从请求头获取 token
	claims, err := jwts.ParseToken(token)
	if err != nil {
		// token 无效
		res.FailWithMsg("认证失败", c)
		c.Abort()
		return
	}

	// 校验角色是否为管理员
	if claims.Role != 1 {
		res.FailWithMsg("无权限访问", c)
		c.Abort()
		return
	}

	// 验证通过，继续执行
	c.Next()
}
