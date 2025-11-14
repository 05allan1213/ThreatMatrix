package log_service

// File: service/log_service/login_log.go
// Description: 用户登录日志记录服务，负责记录用户登录成功/失败日志

import (
	"honey_server/global"
	"honey_server/models"

	"github.com/gin-gonic/gin"
)

// LoginLogService 登录日志服务结构体，用于封装用户 IP 与地址信息
type LoginLogService struct {
	IP   string // 用户 IP 地址
	Addr string // 用户所在地址
}

// NewLoginLog 创建登录日志服务实例
func NewLoginLog(c *gin.Context) *LoginLogService {
	return &LoginLogService{
		IP:   c.ClientIP(), // 获取客户端 IP
		Addr: "",
	}
}

// 记录登录成功日志
func (l LoginLogService) SuccessLog(userID uint, username string) {
	l.save(userID, username, "", "登录成功", true)
}

// 记录登录失败日志
func (l LoginLogService) FailLog(username string, password string, title string) {
	l.save(0, username, password, title, false)
}

// 日志统一写入方法，成功与失败均由此函数落库
func (l LoginLogService) save(userID uint, username string, password string, title string, loginStatus bool) {
	global.DB.Create(&models.LogModel{
		Type:        1,           // 1 代表登录日志
		IP:          l.IP,        // 客户端 IP
		Addr:        l.Addr,      // 地址
		UserID:      userID,      // 用户ID
		Username:    username,    // 用户名
		Pwd:         password,    // 密码（失败时记录）
		LoginStatus: loginStatus, // 登录状态 true=成功 false=失败
		Title:       title,       // 日志标题
	})
}
