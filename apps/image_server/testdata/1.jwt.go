package main

// File:testdata/1.jwt.go
// Description: JWT 生成与解析测试示例

import (
	"fmt"
	"image_server/internal/core"
	"image_server/internal/global"
	"image_server/internal/utils/jwts"
)

func main() {
	global.Config = core.ReadConfig()
	token, _ := jwts.GetToken(jwts.ClaimsUserInfo{
		UserID: 1,
		Role:   1,
	})
	fmt.Println(token)
	claims, err := jwts.ParseToken(token)
	fmt.Println(claims, err)
}
