package api

// File: api/enter.go
// Description: 定义API入口，包含各个子模块的API实例。

import (
	"honey_server/internal/api/captcha_api"
	"honey_server/internal/api/log_api"
	"honey_server/internal/api/node_api"
	"honey_server/internal/api/node_network_api"
	"honey_server/internal/api/user_api"
)

// Api 结构体包含各个子模块的API实例。
type Api struct {
	UserApi        user_api.UserApi
	CaptchaApi     captcha_api.CaptchaApi
	LogApi         log_api.LogApi
	NodeApi        node_api.NodeApi
	NodeNetworkApi node_network_api.NodeNetworkApi
}

var App = Api{}
