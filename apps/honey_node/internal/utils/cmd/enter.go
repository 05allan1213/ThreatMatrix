package cmd

// File: utils/cmd/enter.go
// Description: 提供系统命令执行的封装函数，支持普通命令执行、带路径的命令执行及获取命令输出等功能

import (
	"bytes"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Cmd 执行系统命令（不返回输出内容，仅返回执行结果）
func Cmd(command string) (err error) {
	// 创建Cmd结构体，通过sh -c执行命令（兼容复杂命令语法）
	logrus.Infof("执行命令 %s", command)
	cmd := exec.Command("sh", "-c", command)

	// 设置命令输出缓冲区（用于捕获stdout）
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// 执行命令并返回结果
	err = cmd.Run()
	if err != nil {
		return err
	}

	// 记录命令输出日志
	logrus.Infof("命令输出 %s", stdout.String())
	return nil
}

// Command 执行系统命令并返回输出内容
func Command(command string) (msg string, err error) {
	// 创建Cmd结构体，通过sh -c执行命令
	logrus.Infof("执行命令 %s", command)
	cmd := exec.Command("sh", "-c", command)

	// 设置命令输出缓冲区
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	// 记录命令输出日志并返回输出内容
	logrus.Infof("命令输出 %s", stdout.String())
	return stdout.String(), nil
}

// PathCmd 指定执行路径执行系统命令（不返回输出内容，仅返回执行结果）
func PathCmd(path string, command string) (err error) {
	// 创建Cmd结构体，通过sh -c执行命令
	logrus.Infof("执行命令 %s", command)
	cmd := exec.Command("sh", "-c", command)

	// 设置命令输出缓冲区
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// 指定命令执行的工作路径
	cmd.Path = path

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return err
	}

	// 记录命令输出日志
	logrus.Infof("命令输出 %s", stdout.String())
	return nil
}
