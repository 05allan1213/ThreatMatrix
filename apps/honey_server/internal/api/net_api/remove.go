package net_api

// File: api/net_api/remove.go
// Description: 网络删除API

import (
	"fmt"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// RemoveView 处理网络批量删除的请求
func (NetApi) RemoveView(c *gin.Context) {
	// 获取并绑定批量删除的请求参数（包含待删除的网络ID列表）
	cr := middleware.GetBind[models.IDListRequest](c)
	// 获取请求上下文关联的日志实例
	log := middleware.GetLog(c)

	// 调用通用删除服务执行批量删除操作
	// - 第一个参数：网络模型实例，指定要删除的数据类型
	// - 第二个参数：删除请求配置（包含ID列表、日志实例、操作对象名称"网络"）
	successCount, err := common_service.Remove(models.NetModel{}, common_service.RemoveRequest{
		IDList: cr.IdList, // 待删除的网络ID列表
		Log:    log,       // 日志实例，用于记录删除过程
		Msg:    "网络",      // 操作对象名称，用于日志和提示信息
	})

	// 处理删除操作结果
	if err != nil {
		// 删除失败，返回包含错误信息的提示
		msg := fmt.Sprintf("删除网络失败 %s", err)
		res.FailWithMsg(msg, c)
		return
	}

	// 删除成功，返回包含总数量和成功数量的提示
	msg := fmt.Sprintf("删除成功 共%d个，成功%d个", len(cr.IdList), successCount)
	res.OkWithMsg(msg, c)
}
