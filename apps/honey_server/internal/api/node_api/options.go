package node_api

// File: api/node_api/options.go
// Description: 提供节点选项数据的API接口

import (
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// OptionsResponse 节点选项响应结构体
type OptionsResponse struct {
	Label string `json:"label"` // 节点标题(节点IP)
	Value uint   `json:"value"` // 节点ID
}

// OptionsView 获取节点选项列表的API处理函数
func (NodeApi) OptionsView(c *gin.Context) {
	// 定义节点模型切片，用于存储从数据库查询的节点列表
	var nodeList []models.NodeModel
	// 从数据库查询所有节点记录
	global.DB.Find(&nodeList)

	// 初始化选项列表切片
	var list = make([]OptionsResponse, 0)
	// 遍历节点列表，转换为选项格式
	for _, model := range nodeList {
		list = append(list, OptionsResponse{
			Value: model.ID,                                     // 节点ID
			Label: fmt.Sprintf("%s(%s)", model.Title, model.IP), // 显示标签格式为"节点标题(节点IP)"
		})
	}

	// 返回成功响应，包含转换后的选项列表数据
	res.OkWithData(list, c)
}
