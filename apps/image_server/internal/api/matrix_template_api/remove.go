package matrix_template_api

// File: api/matrix_template_api/remove.go
// Description: 矩阵模板删除

import (
	"fmt"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// 批量删除矩阵模板的API
func (MatrixTemplateApi) Remove(c *gin.Context) {
	// 绑定请求参数，获取要删除的矩阵模板ID列表
	cr := middleware.GetBind[models.IDListRequest](c)
	// 获取上下文日志实例，用于记录删除过程
	log := middleware.GetLog(c)

	// 调用通用删除服务执行删除操作
	successCount, err := common_service.Remove(models.MatrixTemplateModel{}, common_service.RemoveRequest{
		IDList: cr.IdList,
		Log:    log,
		Msg:    "矩阵模板",
	})

	// 处理删除错误
	if err != nil {
		msg := fmt.Sprintf("删除矩阵模板失败 %s", err)
		res.FailWithMsg(msg, c)
		return
	}

	// 返回删除结果（总请求数和成功删除数）
	msg := fmt.Sprintf("删除成功 共%d个，成功%d个", len(cr.IdList), successCount)
	res.OkWithMsg(msg, c)
}
