package main

// File:testdata/2.pwd.go
// Description: 密码哈希生成与验证测试示例

import (
	"fmt"
	"honey_server/internal/utils/pwd"
)

func main() {
	hashPwd, _ := pwd.GenerateFromPassword("1234")
	fmt.Println(hashPwd)
	fmt.Println(pwd.CompareHashAndPassword(hashPwd, "1234"))

}
