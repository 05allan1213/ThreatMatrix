package api

import "image_server/internal/api/mirror_cloud_api"

// File: api/enter.go
// Description: 定义API入口，包含各个子模块的API实例。

// Api 结构体包含各个子模块的API实例。
type Api struct {
	MirrorCloudApi mirror_cloud_api.MirrorCloudApi
}

var App = Api{}
