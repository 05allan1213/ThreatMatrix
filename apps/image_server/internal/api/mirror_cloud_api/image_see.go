package mirror_cloud_api

// File: api/mirror_cloud_api/image_see.go
// Description: 实现镜像预览接口，负责接收上传的 Docker 镜像文件，校验格式和大小，解析镜像元数据，并返回给前端。

import (
	"fmt"
	"image_server/internal/utils/docker"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ImageSeeResponse 镜像信息返回结构体
type ImageSeeResponse struct {
	ImageID   string `json:"imageID"`   // 镜像 ID
	ImageName string `json:"imageName"` // 镜像名称
	ImageTag  string `json:"imageTag"`  // 镜像 Tag
	ImagePath string `json:"imagePath"` // 镜像文件保存路径
}

const (
	maxFileSize  = 2 << 30                // 最大支持文件大小（2GB）
	tempImageDir = "uploads/images_temp/" // 临时保存镜像的目录
)

// ImageSeeView 镜像文件上传与解析接口
func (MirrorCloudApi) ImageSeeView(c *gin.Context) {
	// 获取上传文件
	file, err := c.FormFile("file")
	if err != nil {
		res.FailWithMsg("请选择镜像文件", c)
		return
	}

	// 校验文件大小
	if file.Size > maxFileSize {
		res.FailWithMsg("镜像文件大小不能超过2GB", c)
		return
	}

	// 校验文件格式
	ext := filepath.Ext(file.Filename)
	if ext != ".tar" && ext != ".gz" {
		res.FailWithMsg("只支持.tar和.tar.gz格式的镜像文件", c)
		return
	}

	// 创建临时目录
	if err := os.MkdirAll(tempImageDir, 0755); err != nil {
		res.FailWithMsg(fmt.Sprintf("创建临时目录失败: %v", err), c)
		return
	}

	// 临时保存文件
	tempFilePath := filepath.Join(tempImageDir, file.Filename)
	if err := c.SaveUploadedFile(file, tempFilePath); err != nil {
		res.FailWithMsg(fmt.Sprintf("保存镜像文件失败: %v", err), c)
		return
	}

	// 解析镜像元数据
	imageID, imageName, imageTag, err := docker.ParseImageMetadata(tempFilePath)
	if err != nil {
		os.Remove(tempFilePath)
		res.FailWithMsg(fmt.Sprintf("解析镜像元数据失败: %v", err), c)
		return
	}

	// 异步删除临时文件
	go func() {
		time.Sleep(5 * time.Minute)
		err = os.Remove(tempFilePath)
		if os.IsNotExist(err) {
			return
		}
		if err != nil {
			logrus.Errorf("镜像删除失败 %s", err)
		} else {
			logrus.Infof("删除镜像文件 %s", tempFilePath)
		}
	}()

	// 组织响应数据
	data := ImageSeeResponse{
		ImageID:   imageID,
		ImageName: imageName,
		ImageTag:  imageTag,
		ImagePath: tempFilePath,
	}

	res.OkWithData(data, c)
}
