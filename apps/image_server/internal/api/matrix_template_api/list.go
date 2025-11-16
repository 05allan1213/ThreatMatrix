package matrix_template_api

// File: api/matrix_template_api/list.go
// Description: 矩阵模板列表接口
import (
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ListResponse 矩阵模板列表的响应结构
type ListResponse struct {
	models.MatrixTemplateModel                    // 嵌入矩阵模板基础模型（含ID、标题、创建时间等）
	HostTemplateList           []HostTemplateInfo `json:"hostTemplateList"` // 关联的主机模板列表详情
}

// HostTemplateInfo 矩阵模板中关联的主机模板详情结构
type HostTemplateInfo struct {
	HostTemplateID    uint   `json:"hostTemplateID"`    // 主机模板ID
	HostTemplateTitle string `json:"hostTemplateTitle"` // 主机模板标题
	Weight            int    `json:"weight"`            // 主机模板在矩阵中的权重
}

// ListView 获取矩阵模板列表的API入口函数
//
// 1. 绑定分页参数
// 2. 查询矩阵模板列表（支持标题模糊查询、分页、按创建时间倒序）
// 3. 收集模板关联的所有主机模板ID
// 4. 查询主机模板详情并映射
// 5. 组装包含主机模板详情的响应数据并返回
func (MatrixTemplateApi) ListView(c *gin.Context) {
	// 绑定分页参数（页码、每页条数）
	cr := middleware.GetBind[models.PageInfo](c)

	// 调用通用查询服务获取矩阵模板列表
	_list, count, _ := common_service.QueryList(models.MatrixTemplateModel{},
		common_service.QueryListRequest{
			Likes:    []string{"title"},
			PageInfo: cr,
			Sort:     "created_at desc",
		})

	// 初始化响应列表
	var list = make([]ListResponse, 0)
	// 收集所有矩阵模板关联的主机模板ID，用于批量查询
	var hostTempIDList []uint
	for _, model := range _list {
		for _, port := range model.HostTemplateList {
			hostTempIDList = append(hostTempIDList, port.HostTemplateID)
		}
	}

	// 批量查询关联的主机模板详情
	var hostTemps []models.HostTemplateModel
	global.DB.Find(&hostTemps, "id in ?", hostTempIDList)
	// 将主机模板按ID映射，便于快速查询
	var hostTempMap = map[uint]models.HostTemplateModel{}
	for _, i2 := range hostTemps {
		hostTempMap[i2.ID] = i2
	}

	// 组装响应数据：为每个矩阵模板补充关联的主机模板详情
	for _, model := range _list {
		hostTemplateList := make([]HostTemplateInfo, 0)
		for _, port := range model.HostTemplateList {
			hostTemplateList = append(hostTemplateList, HostTemplateInfo{
				HostTemplateID:    port.HostTemplateID,
				HostTemplateTitle: hostTempMap[port.HostTemplateID].Title, // 从主机模板映射中获取标题
				Weight:            port.Weight,                            // 保留原权重信息
			})
		}
		list = append(list, ListResponse{
			MatrixTemplateModel: model,
			HostTemplateList:    hostTemplateList,
		})
	}

	// 返回列表数据及总条数
	res.OkWithList(list, count, c)
}
