package models

// File: models/log_model.go
// Description: 定义系统日志记录的数据模型及其与用户、服务的关联关系。

// 日志模型
type LogModel struct {
	Model
	Type        int8   `json:"type"`        // 日志类型 1=登录日志
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
