package mirror_cloud_api

// File: api/mirror_cloud_api/image_create.go
// Description: 镜像云端创建接口，实现镜像导入、文件移动、数据库写入等完整流程

import (
	"errors"
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/path"
	"image_server/internal/utils/res"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ImageCreateRequest 请求结构体，用于接收镜像创建所需的各项参数
type ImageCreateRequest struct {
	ImageID   string `json:"imageID" binding:"required"`   // 镜像ID
	ImageName string `json:"imageName" binding:"required"` // 镜像名称
	ImageTag  string `json:"imageTag" binding:"required"`  // 镜像Tag
	ImagePath string `json:"imagePath" binding:"required"` // 镜像上传的路径
	Title     string `json:"title" binding:"required"`     // 镜像别名
	Port      int    `json:"port" binding:"required"`      // 镜像端口
	Agreement int8   `json:"agreement" binding:"required"` // 镜像的协议
}

// ImageCreateView 镜像创建接口，实现镜像导入、校验、移动文件、数据入库等功能
func (MirrorCloudApi) ImageCreateView(c *gin.Context) {
	cr := middleware.GetBind[ImageCreateRequest](c)

	// 1. 检查镜像文件是否存在
	if _, err := os.Stat(cr.ImagePath); errors.Is(err, os.ErrNotExist) {
		res.FailWithMsg("镜像文件不存在", c)
		return
	}

	// 2. 检查镜像title是否重名
	var titleExists models.ImageModel
	if err := global.DB.Take(&titleExists, "title = ?", cr.Title).Error; err == nil {
		res.FailWithMsg("镜像别名不能重复", c)
		return
	}

	// 3. 检查镜像名+tag是否重复
	var nameTagExists models.ImageModel
	if err := global.DB.Take(&nameTagExists, "image_name = ? AND tag = ?", cr.ImageName, cr.ImageTag).Error; err == nil {
		res.FailWithMsg("镜像名称和标签组合不能重复", c)
		return
	}

	// 4. 使用docker load命令导入镜像
	cmd := exec.Command("docker", "load", "-i", cr.ImagePath)
	// 移动到我们的项目路径下
	cmd.Dir = path.GetRootPath()
	output, err := cmd.CombinedOutput()
	if err != nil {
		res.FailWithMsg(fmt.Sprintf("镜像导入失败: %s, 输出: %s", err.Error(), string(output)), c)
		return
	}
	fmt.Println(string(output))

	// 5. 移动镜像文件到正式目录
	finalDir := "uploads/images/"

	// 确保目标目录存在
	if err := os.MkdirAll(finalDir, 0755); err != nil {
		res.FailWithMsg(fmt.Sprintf("创建目标目录失败: %s", err.Error()), c)
		return
	}

	// 获取文件名
	_, fileName := filepath.Split(cr.ImagePath)
	finalPath := filepath.Join(finalDir, fileName)

	// 移动文件
	if err := os.Rename(cr.ImagePath, finalPath); err != nil {
		// 如果移动失败，尝试复制后删除
		logrus.Errorf("文件移动失败 %s", err)
		res.FailWithMsg("文件移动失败", c)
		return
	}

	// 6. 数据入库
	imageModel := models.ImageModel{
		DockerImageID: cr.ImageID,
		ImageName:     cr.ImageName,
		Tag:           cr.ImageTag,
		ImagePath:     finalPath,
		Title:         cr.Title,
		Port:          cr.Port,
		Agreement:     cr.Agreement,
		Status:        1,
	}

	if err := global.DB.Create(&imageModel).Error; err != nil {
		res.FailWithMsg(fmt.Sprintf("数据库插入失败: %s", err.Error()), c)
		return
	}

	res.Ok(imageModel.ID, "镜像创建成功", c)
}
