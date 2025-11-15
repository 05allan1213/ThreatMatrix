package mirror_cloud_api

// File: api/mirror_cloud_api/image_remove.go
// Description: 镜像删除接口，支持检查依赖虚拟服务并删除镜像

import (
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ImageRemoveView 镜像删除接口
// 1. 根据 ID 查询镜像
// 2. 检查是否存在关联虚拟服务
// 3. 执行删除操作
func (MirrorCloudApi) ImageRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	log := middleware.GetLog(c)

	var model models.ImageModel
	err := global.DB.Preload("ServiceList").Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("镜像不存在", c)
		return
	}

	// 检查是否存在虚拟服务依赖
	if len(model.ServiceList) > 0 {
		res.FailWithMsg("镜像存在虚拟服务，请先删除关联的虚拟服务", c)
		return
	}

	log.Infof("删除镜像 %#v", model)

	// 执行删除
	err = global.DB.Delete(&model).Error
	if err != nil {
		log.Errorf("删除镜像失败 %s", err)
		res.FailWithMsg("镜像删除失败", c)
		return
	}

	res.OkWithMsg("删除成功", c)
}
