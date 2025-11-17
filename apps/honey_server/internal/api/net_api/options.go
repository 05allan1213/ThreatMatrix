package net_api

// File: api/net_api/options.go
// Description: 网络选项列表API

import (
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// OptionsResponse 网络选项响应结构体
type OptionsResponse struct {
	Label string `json:"label"` // 显示标签，格式为"网络标题(子网信息)"
	Value uint   `json:"value"` // 网络ID，作为选项对应的值
}

// OptionsView 处理获取网络选项列表的请求
func (NetApi) OptionsView(c *gin.Context) {
	// 查询所有网络列表数据
	var netList []models.NetModel
	global.DB.Find(&netList)

	// 转换网络列表为选项响应格式
	var list = make([]OptionsResponse, 0)
	for _, model := range netList {
		list = append(list, OptionsResponse{
			Value: model.ID,                                           // 网络ID作为选项值
			Label: fmt.Sprintf("%s(%s)", model.Title, model.Subnet()), // 标签格式：网络标题(子网信息)
		})
	}

	// 返回转换后的选项列表数据
	res.OkWithData(list, c)
}
