package cron_service

// File: service/cron_service/resource.go
// Description: 初始化并启动定时任务调度器

import (
	"time"

	"github.com/robfig/cron/v3"
)

// 初始化并启动定时任务调度器
func Run() {
	// 加载上海时区（Asia/Shanghai），用于定时任务的时间计算
	timezone, _ := time.LoadLocation("Asia/Shanghai")

	// 创建定时任务调度器，配置支持秒级精度（WithSeconds()）和指定时区（WithLocation(timezone)）
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))

	// 向调度器添加任务：每5秒执行一次Resource函数（节点资源信息上报）
	crontab.AddFunc("*/5 * * * * *", Resource)

	// 启动定时任务调度器，开始执行已添加的任务
	crontab.Start()
}
