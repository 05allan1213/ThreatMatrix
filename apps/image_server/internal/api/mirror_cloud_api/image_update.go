package mirror_cloud_api

// File: api/mirror_cloud_api/image_update.go
// Description: 镜像更新接口，支持基础信息修改与重名校验

import (
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ImageUpdateRequest 镜像更新请求结构体
// 包含镜像名称、端口、协议、状态、Logo、描述等可修改字段
type ImageUpdateRequest struct {
	ID        uint   `json:"id"`                                      // 镜像 ID
	Title     string `json:"title" binding:"required"`                // 镜像别名
	Port      int    `json:"port" binding:"required,min=1,max=65535"` // 镜像端口
	Agreement int8   `json:"agreement" binding:"required,oneof=1"`    // 镜像协议
	Status    int8   `json:"status" binding:"required,oneof=1 2"`     // 镜像状态 1 成功  2 禁用
	Logo      string `json:"logo"`                                    // 镜像 Logo
	Desc      string `json:"desc"`                                    // 镜像描述
}

// ImageUpdateView 镜像更新接口
// 1. 校验镜像是否存在
// 2. 校验 title 是否重名
// 3. 执行更新操作
func (MirrorCloudApi) ImageUpdateView(c *gin.Context) {
	cr := middleware.GetBind[ImageUpdateRequest](c)

	var model models.ImageModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("镜像不存在", c)
		return
	}

	// title 不能和其他镜像重名
	var newModel models.ImageModel
	err = global.DB.Take(&newModel, "id <> ? and title = ?", cr.ID, cr.Title).Error
	if err == nil {
		res.FailWithMsg("修改的镜像名称不能重复", c)
		return
	}

	err = global.DB.Model(&model).Updates(models.ImageModel{
		Title:     cr.Title,
		Port:      cr.Port,
		Agreement: cr.Agreement,
		Status:    cr.Status,
		Logo:      cr.Logo,
		Desc:      cr.Desc,
	}).Error
	if err != nil {
		res.FailWithMsg("镜像更新失败", c)
		return
	}

	res.OkWithMsg("镜像修改成功", c)
}
