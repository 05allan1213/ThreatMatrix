package core

// File: core/logger.go
// Description: 定义日志系统的初始化、自定义格式化输出以及日志文件轮转逻辑。

import (
	"bytes"
	"fmt"
	"honey_server/global"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// MyLog 自定义日志格式化器（实现 logrus.Formatter 接口）
// 用于控制日志在控制台输出的样式（颜色、时间、文件、函数名等）
type MyLog struct{}

// 定义颜色常量（ANSI 颜色码，用于终端高亮输出）
const (
	red    = 31 // 错误、严重日志
	yellow = 33 // 警告日志
	blue   = 36 // 普通信息日志
	gray   = 37 // 调试日志
)

// Format 自定义日志输出格式（控制控制台输出的颜色、时间、调用位置等）
func (MyLog) Format(entry *logrus.Entry) ([]byte, error) {
	// 根据不同的日志级别设置颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	// 复用缓冲区以提升性能
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 自定义时间格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")

	// 如果开启了 ReportCaller，则输出调用文件与函数信息
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)

		// 自定义控制台输出格式：
		// appName 时间 [日志级别] 文件:行号 函数名 日志内容
		fmt.Fprintf(b, "%s [%s] \x1b[%dm[%s]\x1b[0m %s %s %s\n",
			entry.Data["appName"], timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message)
	}

	return b.Bytes(), nil
}

// MyHook 自定义 Hook，用于日志文件输出与按日期轮转
type MyHook struct {
	file     *os.File // info.log 文件句柄
	errFile  *os.File // err.log 文件句柄
	fileDate string   // 当前日志文件日期（用于判断是否需要轮转）
	logPath  string   // 日志根目录
	mu       sync.Mutex
}

// Fire 是 Hook 的核心逻辑：日志触发时写入文件
func (hook *MyHook) Fire(entry *logrus.Entry) error {
	hook.mu.Lock()
	defer hook.mu.Unlock()

	// 当前日期（用于按天轮转日志）
	timer := entry.Time.Format("2006-01-02")

	// 将日志条目格式化为字符串
	line, err := entry.String()
	if err != nil {
		return fmt.Errorf("格式化日志内容失败: %v", err)
	}

	// 如果日期发生变化，则执行日志轮转
	if hook.fileDate != timer {
		if err := hook.rotateFiles(timer); err != nil {
			return err
		}
	}

	// 写入 info.log
	if _, err := hook.file.Write([]byte(line)); err != nil {
		return fmt.Errorf("写入 info.log 失败: %v", err)
	}

	// 如果是错误级别或更高，则额外写入 err.log
	if entry.Level <= logrus.ErrorLevel {
		if _, err := hook.errFile.Write([]byte(line)); err != nil {
			return fmt.Errorf("写入 err.log 失败: %v", err)
		}
	}

	return nil
}

// rotateFiles 日志轮转：按日期创建新的日志目录与文件
func (hook *MyHook) rotateFiles(timer string) error {
	// 关闭旧文件句柄
	if hook.file != nil {
		if err := hook.file.Close(); err != nil {
			return fmt.Errorf("关闭 info.log 失败: %v", err)
		}
	}
	if hook.errFile != nil {
		if err := hook.errFile.Close(); err != nil {
			return fmt.Errorf("关闭 err.log 失败: %v", err)
		}
	}

	// 按日期创建目录，例如 logs/2025-11-10/
	dirName := fmt.Sprintf("%s/%s", hook.logPath, timer)
	if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 构造文件路径
	infoFilename := fmt.Sprintf("%s/info.log", dirName)
	errFilename := fmt.Sprintf("%s/err.log", dirName)

	// 打开 info.log
	var err error
	hook.file, err = os.OpenFile(infoFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("打开 info.log 失败: %v", err)
	}

	// 打开 err.log
	hook.errFile, err = os.OpenFile(errFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("打开 err.log 失败: %v", err)
	}

	hook.fileDate = timer
	return nil
}

// Levels 指定哪些日志级别会触发 Hook（这里选择所有级别）
func (hook *MyHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// GetLogger 初始化日志系统，返回带有应用名字段的 *logrus.Entry
func GetLogger() *logrus.Entry {
	logger := logrus.New()
	l := global.Config.Logger

	// 解析日志级别
	level, err := logrus.ParseLevel(l.Level)
	if err != nil {
		logrus.Warnf("日志级别配置错误，已自动切换为 info")
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 添加自定义 Hook（负责文件写入与轮转）
	logger.AddHook(&MyHook{logPath: "logs"})

	// 根据配置选择日志输出格式
	if l.Format == "json" {
		// JSON 格式（适合生产环境或日志收集系统）
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.DateTime,
		})
	} else {
		// 自定义终端输出格式（带颜色、时间、调用信息）
		logger.SetFormatter(&MyLog{})
	}

	// 启用调用者追踪（用于输出函数名与行号）
	logger.SetReportCaller(true)

	// 返回带有应用名字段的日志对象
	return logger.WithField("appName", l.AppName)
}
