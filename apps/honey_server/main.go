package main

import (
	"honey_server/core"
	"honey_server/flags"
	"honey_server/global"
)

func main() {
	global.DB = core.InitDB()
	flags.Run()
}
