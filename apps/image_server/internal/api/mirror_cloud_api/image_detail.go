package mirror_cloud_api

// File: api.mirror_cloud_api/image_detail.go
// Description: 镜像详情接口，通过 ID 查询镜像完整信息

import (
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ImageDetailView 镜像详情查询接口
// 根据传入的 IDRequest 查询对应的镜像记录
func (MirrorCloudApi) ImageDetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var model models.ImageModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("镜像不存在", c)
		return
	}

	res.OkWithData(model, c)
}
