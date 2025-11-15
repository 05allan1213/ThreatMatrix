package host_template_api

// File: api/host_template_api/list.go
// Description: 主机模板列表 API，负责查询主机模板列表，支持分页和标题模糊查询，返回模板信息及关联的虚拟服务详情（包括服务标题、状态等）

import (
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/service/common_service"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// ListResponse 主机模板列表的响应结构
type ListResponse struct {
	models.HostTemplateModel                        // 嵌入主机模板基础模型（含ID、标题、创建时间等）
	PortList                 []HostTemplatePortInfo `json:"portList"` // 端口列表详情（含关联的虚拟服务信息）
}

// HostTemplatePortInfo 主机模板中端口关联的虚拟服务详情结构
type HostTemplatePortInfo struct {
	Port          int    `json:"port"`          // 端口号
	ServiceID     uint   `json:"serviceID"`     // 关联的虚拟服务ID
	ServiceTitle  string `json:"serviceTitle"`  // 虚拟服务标题
	ServiceStatus int8   `json:"serviceStatus"` // 虚拟服务状态（1：正常，2：异常）
}

// ListView 获取主机模板列表的API入口函数
//
// 1. 绑定分页参数
// 2. 查询主机模板列表（支持标题模糊查询、分页、按创建时间倒序）
// 3. 收集模板关联的所有虚拟服务ID
// 4. 查询虚拟服务详情并映射
// 5. 组装包含服务详情的响应数据并返回
func (HostTemplateApi) ListView(c *gin.Context) {
	// 绑定分页参数（页码、每页条数）
	cr := middleware.GetBind[models.PageInfo](c)

	// 调用通用查询服务获取主机模板列表
	_list, count, _ := common_service.QueryList(models.HostTemplateModel{},
		common_service.QueryListRequest{
			Likes:    []string{"title"},
			PageInfo: cr,
			Sort:     "created_at desc",
		})

	// 初始化响应列表
	var list = make([]ListResponse, 0)
	// 收集所有模板关联的虚拟服务ID，用于批量查询
	var serviceIDList []uint
	for _, model := range _list {
		for _, port := range model.PortList {
			serviceIDList = append(serviceIDList, port.ServiceID)
		}
	}

	// 批量查询关联的虚拟服务详情
	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList, "id in ?", serviceIDList)
	// 将服务按ID映射，便于快速查询
	var serviceMap = map[uint]models.ServiceModel{}
	for _, i2 := range serviceList {
		serviceMap[i2.ID] = i2
	}

	// 组装响应数据：为每个模板补充端口关联的服务详情
	for _, model := range _list {
		portList := make([]HostTemplatePortInfo, 0)
		for _, port := range model.PortList {
			portList = append(portList, HostTemplatePortInfo{
				Port:          port.Port,
				ServiceID:     port.ServiceID,
				ServiceTitle:  serviceMap[port.ServiceID].Title,  // 从服务映射中获取标题
				ServiceStatus: serviceMap[port.ServiceID].Status, // 从服务映射中获取状态
			})
		}
		list = append(list, ListResponse{
			HostTemplateModel: model,
			PortList:          portList,
		})
	}

	// 返回列表数据及总条数
	res.OkWithList(list, count, c)
}
