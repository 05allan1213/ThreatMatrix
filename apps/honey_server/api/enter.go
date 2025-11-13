package api

// File: api/enter.go
// Description: 定义API入口，包含各个子模块的API实例。

import (
	"honey_server/api/captcha_api"
	"honey_server/api/user_api"
)

// Api 结构体包含各个子模块的API实例。
type Api struct {
	UserApi    user_api.UserApi
	CaptchaApi captcha_api.CaptchaApi
}

var App = Api{}
