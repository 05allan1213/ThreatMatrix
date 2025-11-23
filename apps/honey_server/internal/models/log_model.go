package models

// File: models/log_model.go
// Description: 定义系统日志记录的数据模型及其与用户、服务的关联关系。

// 日志模型
type LogModel struct {
	Model
	Type        int8   `json:"type"`                            // 日志类型 1=登录日志
	IP          string `gorm:"size:32;index:idx_ip" json:"ip"`  // IP地址
	Addr        string `gorm:"size:64" json:"addr"`             // 地址
	UserID      uint   `gorm:"index:idx_user_id" json:"userID"` // 用户ID
	Username    string `gorm:"size:32" json:"username"`         // 用户名
	Pwd         string `gorm:"size:64" json:"pwd"`              // 密码
	LoginStatus bool   `json:"loginStatus"`                     // 登录状态
	Title       string `gorm:"size:64" json:"title"`            // 日志标题
	Level       int8   `json:"level"`                           // 日志级别
	Content     string `json:"content"`                         // 日志内容
	ServiceName string `gorm:"size:32" json:"serviceName"`      // 服务名称
}
