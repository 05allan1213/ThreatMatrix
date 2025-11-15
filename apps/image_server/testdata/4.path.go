package main

// File: /testdata/4.path.go
// Description: 获取项目根路径

import (
	"fmt"
	"image_server/internal/utils/path"
)

func main() {
	fmt.Println(path.GetRootPath())
}
