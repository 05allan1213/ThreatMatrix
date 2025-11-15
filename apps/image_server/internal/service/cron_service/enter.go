package cron_service

// File: service/cron_service/enter.go
// Description: 定时任务调度器，负责启动所有 cron 任务。

import (
	"time"

	"github.com/robfig/cron/v3"
)

// Run 启动定时任务调度器
// 1. 设置时区为上海时间
// 2. 使用 cron.WithSeconds() 允许秒级任务
// 3. 注册容器健康检查任务（每分钟执行一次）
func Run() {
	// 加载亚洲/上海时区，作为定时任务的时间基准
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	// 创建cron实例，配置支持秒级定时任务，并使用上海时区
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	// 注册定时任务：每分钟执行一次VsHealth函数
	crontab.AddFunc("0 * * * * *", VsHealth)
	// 启动定时任务调度器，开始执行已注册的任务
	crontab.Start()
}
