package info

// File: utils/info/system.go
// Description: 系统基础信息获取

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// SystemInfo 系统基础信息结构体
type SystemInfo struct {
	OSVersion    string // 操作系统发行版本
	Kernel       string // 内核版本
	Architecture string // 系统架构
	BootTime     string // 系统启动时间
}

// 获取系统基础信息
func GetSystemInfo() (SystemInfo, error) {
	info := SystemInfo{}
	var err error

	// 获取操作系统发行版本
	info.OSVersion, err = getOSVersion()
	if err != nil {
		return info, fmt.Errorf("获取发行版本失败: %v", err)
	}

	// 获取内核版本
	info.Kernel, err = getKernelVersion()
	if err != nil {
		return info, fmt.Errorf("获取内核版本失败: %v", err)
	}

	// 获取系统架构（如x86_64、arm64等）
	info.Architecture, err = getArchitecture()
	if err != nil {
		return info, fmt.Errorf("获取系统架构失败: %v", err)
	}

	// 获取系统启动时间（格式化后）
	info.BootTime, err = getBootTime()
	if err != nil {
		return info, fmt.Errorf("获取系统启动时间失败: %v", err)
	}

	return info, nil
}

// 读取/etc/os-release文件获取操作系统发行版本
func getOSVersion() (string, error) {
	// 打开系统发行版本信息文件（Linux系统标准路径）
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer file.Close() // 确保文件最终关闭

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 查找包含PRETTY_NAME的行（该字段通常存储可读的发行版本名称）
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			// 按"="分割，取后半部分（发行版本名称）
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				// 去除名称前后的引号（单引号或双引号）
				name := strings.Trim(parts[1], "\"'")
				return name, nil
			}
		}
	}

	// 若未找到PRETTY_NAME字段，返回错误
	return "", fmt.Errorf("未找到发行版本信息")
}

// 执行uname -r命令获取内核版本
func getKernelVersion() (string, error) {
	// 执行uname -r命令（该命令用于显示内核版本）
	cmd := exec.Command("uname", "-r")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// 转换输出为字符串并去除首尾空白（如换行符）
	return strings.TrimSpace(string(output)), nil
}

// 执行uname -m命令获取系统架构
func getArchitecture() (string, error) {
	// 执行uname -m命令（该命令用于显示机器硬件名称，即架构）
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	// 转换输出为字符串并去除首尾空白
	return strings.TrimSpace(string(output)), nil
}

// 读取/proc/stat文件获取系统启动时间并格式化
func getBootTime() (string, error) {
	// 打开系统状态信息文件（Linux系统标准路径，包含系统启动时间等信息）
	file, err := os.Open("/proc/stat")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 查找包含btime的行（该字段存储系统启动的Unix时间戳）
		if strings.HasPrefix(line, "btime ") {
			// 按空白分割，取第二个元素（时间戳）
			parts := strings.Fields(line)
			if len(parts) == 2 {
				// 将时间戳字符串转换为int64类型
				btime, err := strconv.ParseInt(parts[1], 10, 64)
				if err == nil {
					// 将Unix时间戳转换为"年-月-日 时:分:秒"格式
					bootTime := time.Unix(btime, 0)
					return bootTime.Format("2006-01-02 15:04:05"), nil
				}
			}
		}
	}

	// 若未找到btime字段，返回错误
	return "", fmt.Errorf("未找到系统启动时间信息")
}
