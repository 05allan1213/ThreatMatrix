package host_template_api

// File: api/host_template_api/update.go
// Description: 主机模板更新

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// UpdateRequest 更新主机模板的请求参数结构
type UpdateRequest struct {
	ID       uint                        `json:"id" binding:"required"`    // 主机模板ID
	Title    string                      `json:"title" binding:"required"` // 更新后的模板名称
	PortList models.HostTemplatePortList `json:"portList" binding:"dive"`  // 更新后的端口列表，需包含关联的虚拟服务ID和端口信息
}

// UpdateView 更新主机模板的API入口函数
//
// 1. 绑定并验证请求参数
// 2. 检查待更新的模板是否存在
// 3. 校验新名称是否与其他模板重复
// 4. 校验端口列表中端口是否重复
// 5. 校验关联的虚拟服务是否存在
// 6. 执行更新并返回结果
func (HostTemplateApi) UpdateView(c *gin.Context) {
	// 绑定并验证请求参数（从上下文获取UpdateRequest结构体）
	cr := middleware.GetBind[UpdateRequest](c)

	// 检查待更新的主机模板是否存在
	var model models.HostTemplateModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("主机模板不存在", c)
		return
	}

	// 校验新模板名称是否与其他模板重复（排除当前模板自身）
	var newModel models.HostTemplateModel
	err = global.DB.Take(&newModel, "title = ? and id <> ?", cr.Title, cr.ID).Error
	if err == nil {
		res.FailWithMsg("修改的主机模板名称不能重复", c)
		return
	}

	// 收集端口列表中所有关联的虚拟服务ID，用于后续存在性校验
	var serviceIDList []uint
	// 用于校验端口是否重复的映射表（键为端口号，值为存在标识）
	var portMap = map[int]bool{}
	for _, port := range cr.PortList {
		serviceIDList = append(serviceIDList, port.ServiceID)
		portMap[port.Port] = true // 利用map键唯一性记录端口，用于检测重复
	}

	// 校验端口列表中是否存在重复端口（map长度与原列表长度不一致则说明有重复）
	if len(portMap) != len(cr.PortList) {
		res.FailWithMsg("端口存在重复", c)
		return
	}

	// 查询所有关联的虚拟服务，验证是否存在
	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList, "id in ?", serviceIDList)
	// 将查询到的服务按ID映射，便于快速查找
	var serviceMap = map[uint]models.ServiceModel{}
	for _, serviceModel := range serviceList {
		serviceMap[serviceModel.ID] = serviceModel
	}

	// 校验端口列表中关联的虚拟服务是否均存在
	for _, port := range cr.PortList {
		_, ok := serviceMap[port.ServiceID]
		if !ok {
			msg := fmt.Sprintf("虚拟服务 %d 不存在", port.ServiceID)
			res.FailWithMsg(msg, c)
			return
		}
	}

	// 构建更新数据并执行数据库更新操作
	newModel = models.HostTemplateModel{
		Title:    cr.Title,
		PortList: cr.PortList,
	}
	err = global.DB.Model(&model).Updates(newModel).Error
	if err != nil {
		res.FailWithMsg("主机模板更新失败", c)
		return
	}

	// 返回更新成功消息
	res.OkWithMsg("主机模板更新成功", c)
}
