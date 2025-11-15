package host_template_api

// File: api/host_template_api/options.go
// Description: 主机模板选项列表
import (
	"image_server/internal/global"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// OptionsListResponse 主机模板选项列表的响应结构
type OptionsListResponse struct {
	Label string `json:"label"` // 选项显示标签（主机模板标题）
	Value uint   `json:"value"` // 选项对应的值（主机模板ID）
}

// OptionsView 获取主机模板选项列表的API入口函数
func (HostTemplateApi) OptionsView(c *gin.Context) {
	// 初始化选项列表
	var list = make([]OptionsListResponse, 0)

	// 从数据库查询主机模板的ID和标题，直接映射为选项格式
	global.DB.Model(models.HostTemplateModel{}).Select("id as value", "title as label").Scan(&list)

	// 返回选项列表数据
	res.OkWithData(list, c)
}
