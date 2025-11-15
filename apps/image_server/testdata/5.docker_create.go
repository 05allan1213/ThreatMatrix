package main

// File: /testdata/5.docker_create.go
// Description: 通过 docker 命令运行一个容器，并打印执行耗时

import (
	"fmt"
	"image_server/internal/utils/cmd"
	"time"
)

func main() {
	t1 := time.Now() // 记录开始时间

	// 确保Docker网络存在
	networkCommand := "docker network create --driver bridge --subnet 10.2.0.0/24 honey-hy >/dev/null 2>&1 || true"
	err := cmd.Cmd(networkCommand)
	if err != nil {
		fmt.Println("检查或创建Docker网络失败:", err)
		fmt.Println(time.Since(t1)) // 输出耗时
		return
	}

	// 使用 docker 命令运行容器
	// 示例:
	// docker network create --driver bridge --subnet 10.2.0.0/24 honey-hy
	// docker run -d --network honey-hy --ip 10.2.0.10 --name my_container image_name:tag
	ip := "10.2.0.2"
	command := fmt.Sprintf(
		"docker run -d --network honey-hy --ip %s --name %s_%d %s:%s tail -f /dev/null",
		ip, "alpine", time.Now().Unix(), "alpine", "latest",
	)

	// 执行命令
	err = cmd.Cmd(command)
	if err != nil {
		fmt.Println(err)
		fmt.Println(time.Since(t1)) // 输出耗时
		return
	}

	fmt.Println(time.Since(t1)) // 输出执行总耗时
}