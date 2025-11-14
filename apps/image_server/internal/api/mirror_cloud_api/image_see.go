package mirror_cloud_api

// File: api/mirror_cloud_api/image_see.go
// Description: 镜像云相关接口，提供镜像文件预检查/查看功能

import (
	"fmt"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ImageSeeView 镜像文件查看/检测接口
// 从表单中读取上传的文件并进行基础处理
func (MirrorCloudApi) ImageSeeView(c *gin.Context) {
	file, err := c.FormFile("file") // 获取上传的镜像文件
	if err != nil {
		res.FailWithMsg("请选择镜像文件", c)
		return
	}

	// 打印文件名
	fmt.Println(file.Filename)

	res.OkWithData(nil, c)
}
