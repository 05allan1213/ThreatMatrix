package vs_api

// File: api/vs_api/vs_options.go
// Description: 虚拟服务选项列表 API，用于获取虚拟服务的选项数据

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// VsOptionsListResponse 虚拟服务选项列表的响应结构
// 用于前端下拉选择组件，包含显示文本、对应值及禁用状态
type VsOptionsListResponse struct {
	Label   string `json:"label"`   // 选项显示标签，格式为"服务标题/端口"
	Value   uint   `json:"value"`   // 选项对应的值，即虚拟服务的ID
	Disable bool   `json:"disable"` // 选项是否禁用
}

// VsOptionsListView 获取虚拟服务选项列表的API入口函数
// 功能：查询所有虚拟服务记录，转换为前端下拉组件所需的选项格式并返回
func (VsApi) VsOptionsListView(c *gin.Context) {
	// 查询数据库中所有虚拟服务记录
	var list []models.ServiceModel
	global.DB.Find(&list)

	// 将服务记录转换为选项列表格式
	var options []VsOptionsListResponse
	for _, model := range list {
		// 构建单个选项：标签为"服务标题/端口"，值为服务ID
		item := VsOptionsListResponse{
			Label: fmt.Sprintf("%s/%d", model.Title, model.Port),
			Value: model.ID,
		}
		// 若服务状态非1（非正常状态），则标记为禁用
		if model.Status != 1 {
			item.Disable = true
		}
		options = append(options, item)
	}

	// 返回转换后的选项列表数据
	res.OkWithData(options, c)
}
