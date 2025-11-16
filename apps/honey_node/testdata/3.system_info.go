package main

// File: testdata/3.system_info.go
// Description: 获取系统信息

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// SystemInfo 包含系统信息的结构体
type SystemInfo struct {
	OSVersion    string
	Kernel       string
	Architecture string
	BootTime     string
}

// GetSystemInfo 获取系统信息
func GetSystemInfo() (SystemInfo, error) {
	info := SystemInfo{}
	var err error

	// 获取发行版本
	info.OSVersion, err = getOSVersion()
	if err != nil {
		return info, fmt.Errorf("获取发行版本失败: %v", err)
	}

	// 获取内核版本
	info.Kernel, err = getKernelVersion()
	if err != nil {
		return info, fmt.Errorf("获取内核版本失败: %v", err)
	}

	// 获取系统架构
	info.Architecture, err = getArchitecture()
	if err != nil {
		return info, fmt.Errorf("获取系统架构失败: %v", err)
	}

	// 获取系统启动时间
	info.BootTime, err = getBootTime()
	if err != nil {
		return info, fmt.Errorf("获取系统启动时间失败: %v", err)
	}

	return info, nil
}

// 获取发行版本信息
func getOSVersion() (string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				// 去除引号
				name := strings.Trim(parts[1], "\"'")
				return name, nil
			}
		}
	}

	return "", fmt.Errorf("未找到发行版本信息")
}

// 获取内核版本
func getKernelVersion() (string, error) {
	cmd := exec.Command("uname", "-r")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// 获取系统架构
func getArchitecture() (string, error) {
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// 获取系统启动时间
func getBootTime() (string, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "btime ") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				if btime, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
					bootTime := time.Unix(btime, 0)
					return bootTime.Format("2006-01-02 15:04:05"), nil
				}
			}
		}
	}

	return "", fmt.Errorf("未找到系统启动时间信息")
}

func main() {
	info, err := GetSystemInfo()
	if err != nil {
		fmt.Println("错误:", err)
		return
	}

	fmt.Printf("发行版本: %s\n", info.OSVersion)
	fmt.Printf("内核版本: %s\n", info.Kernel)
	fmt.Printf("系统架构: %s\n", info.Architecture)
	fmt.Printf("系统启动时间: %s\n", info.BootTime)
}
