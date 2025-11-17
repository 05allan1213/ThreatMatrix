package net_api

// File: api/net_api/detail.go
// Description: 网络详情API

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// DetailView 处理获取网络详情的请求
func (NetApi) DetailView(c *gin.Context) {
	// 从请求中获取绑定的ID参数（用于指定要查询的网络）
	cr := middleware.GetBind[models.IDRequest](c)

	// 查询指定ID的网络详情
	var model models.NetModel
	err := global.DB.Take(&model, cr.Id).Error
	if err != nil {
		// 若查询失败（网络不存在），返回错误提示
		res.FailWithMsg("网络不存在", c)
		return
	}

	// 查询成功，返回网络详情数据
	res.OkWithData(model, c)
}
