package api

import (
	"image_server/internal/api/mirror_cloud_api"
	"image_server/internal/api/vs_api"
	"image_server/internal/api/vs_net_api"
)

// File: api/enter.go
// Description: 定义API入口，包含各个子模块的API实例。

// Api 结构体包含各个子模块的API实例。
type Api struct {
	MirrorCloudApi mirror_cloud_api.MirrorCloudApi
	VsApi          vs_api.VsApi
	VsNetApi       vs_net_api.VsNetApi
}

var App = Api{}
