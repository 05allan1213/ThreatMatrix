package middleware

// File: middleware/log_middleware.go
// Description: 日志中间件，为每个请求分配唯一日志ID，方便追踪请求链路。

import (
	"honey_server/global"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// LogMiddleware 日志中间件
func LogMiddleware(c *gin.Context) {
	log := global.Log
	uid := uuid.New().String()

	// 为当前请求创建带 logID 的 logger
	logger := log.WithField("logID", uid)

	// 将 logger 存入 gin.Context，供后续处理函数使用
	c.Set("log", logger)
}

// GetLog 从 gin.Context 获取 logger
// 如果中间件未设置，则会 panic
func GetLog(c *gin.Context) *logrus.Entry {
	return c.MustGet("log").(*logrus.Entry)
}
