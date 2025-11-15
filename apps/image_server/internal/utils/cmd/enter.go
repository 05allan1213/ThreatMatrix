package cmd

// File: utils/cmd/enter.go
// Description: 命令执行工具，支持在默认或指定路径下执行 shell 命令，并记录日志

import (
	"bytes"
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Cmd 在默认 shell 环境下执行命令
func Cmd(command string) (err error) {
	// 创建一个 Cmd 结构体
	logrus.Infof("执行命令 %s", command)
	cmd := exec.Command("sh", "-c", command)

	// 设置输出缓冲
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return err
	}

	logrus.Infof("命令输出 %s", stdout.String())
	return nil
}

// PathCmd 在指定路径下执行命令
func PathCmd(path string, command string) (err error) {
	// 创建一个 Cmd 结构体
	logrus.Infof("执行命令 %s", command)
	cmd := exec.Command("sh", "-c", command)

	// 设置输出缓冲
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	// 设置命令执行路径
	cmd.Path = path

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return err
	}

	logrus.Infof("命令输出 %s", stdout.String())
	return nil
}
