package vs_api

// File: api/vs_api/vs_remove.go
// Description: 虚拟服务删除 API，负责删除指定的虚拟服务记录。

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// VsRemoveView 删除虚拟服务
func (VsApi) VsRemoveView(c *gin.Context) {
	// 绑定请求参数，获取要删除的虚拟服务ID列表
	cr := middleware.GetBind[models.IDListRequest](c)

	// 查询数据库中指定ID列表的虚拟服务记录
	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList, "id in ?", cr.IdList)

	// 检查是否存在对应的虚拟服务记录
	if len(serviceList) == 0 {
		res.FailWithMsg("不存在的虚拟服务", c)
		return
	}

	// 执行删除操作
	result := global.DB.Delete(&serviceList)
	successCount := result.RowsAffected // 记录成功删除的数量
	err := result.Error

	// 处理删除错误
	if err != nil {
		res.FailWithMsg("删除虚拟服务失败", c)
		return
	}

	// 返回删除成功的消息（包含成功删除的数量）
	msg := fmt.Sprintf("删除虚拟服务成功 共%d个", successCount)
	res.OkWithMsg(msg, c)
}
