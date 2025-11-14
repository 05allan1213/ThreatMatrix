package core

// File: core/config.go
// Description: 定义配置文件的读取和解析逻辑。

import (
	"honey_server/internal/config"
	"honey_server/internal/flags"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// 读取配置文件
func ReadConfig() *config.Config {
	byteData, err := os.ReadFile(flags.Options.File)
	if err != nil {
		logrus.Fatalf("配置文件读取错误 %s", err)
		return nil
	}
	var c = new(config.Config)
	err = yaml.Unmarshal(byteData, &c)
	if err != nil {
		logrus.Fatalf("配置文件配置错误 %s", err)
		return nil
	}
	return c
}
