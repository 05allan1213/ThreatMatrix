package main

// File: testdata/delay.go
// Description: 程序功能：打印初始时间，随后分别等待指定的时长（2秒、10秒、20秒），每次等待结束后打印当前时间

import (
	"fmt"
	"time"
)

func main() {
	// 定义需要等待的时长列表（单位：秒）
	delays := []time.Duration{
		2 * time.Second,
		10 * time.Second,
		20 * time.Second,
	}

	// 打印当前时间（初始时间），格式为日期时间
	fmt.Println(time.Now().Format(time.DateTime))

	// 遍历每个等待时长，执行等待后打印当前时间
	for _, d := range delays {
		time.Sleep(d)                                 // 等待指定时长
		fmt.Println(time.Now().Format(time.DateTime)) // 打印等待结束后的当前时间
	}
}
