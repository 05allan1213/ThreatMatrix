package mirror_cloud_api

// File: api/mirror_cloud_api/image_options.go
// Description: 镜像选项列表接口，用于前端下拉选择，可标记禁用状态

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ImageOptionsListResponse 镜像选项列表响应结构体
type ImageOptionsListResponse struct {
	Label   string `json:"label"`   // 前端显示文本
	Value   uint   `json:"value"`   // 镜像 ID
	Disable bool   `json:"disable"` // 是否禁用
}

// ImageOptionsListView 镜像选项列表接口
// 查询所有镜像，将状态为 2 的镜像标记为禁用
func (MirrorCloudApi) ImageOptionsListView(c *gin.Context) {
	var list []models.ImageModel
	global.DB.Find(&list)

	var options []ImageOptionsListResponse
	for _, model := range list {
		item := ImageOptionsListResponse{
			Label: fmt.Sprintf("%s/%d", model.Title, model.Port),
			Value: model.ID,
		}
		if model.Status == 2 {
			item.Disable = true
		}
		options = append(options, item)
	}

	res.OkWithData(options, c)
}
