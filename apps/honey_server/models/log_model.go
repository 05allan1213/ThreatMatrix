// Package models 定义诱捕服务所使用的数据实体。
//
// 本文件描述系统日志记录的结构体，用于保存审计与访问信息。
package models

import "gorm.io/gorm"

// 日志模型
type LogModel struct {
	gorm.Model
	Type        int8   `json:"type"`        // 日志类型
	IP          string `json:"ip"`          // IP地址
	Addr        string `json:"addr"`        // 地址
	UserID      uint   `json:"userID"`      // 用户ID
	Username    string `json:"username"`    // 用户名
	Pwd         string `json:"pwd"`         // 密码
	LoginStatus bool   `json:"loginStatus"` // 登录状态
	Title       string `json:"title"`       // 日志标题
	Level       int8   `json:"level"`       // 日志级别
	Content     string `json:"content"`     // 日志内容
	ServiceName string `json:"serviceName"` // 服务名称
}
