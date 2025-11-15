package host_template_api

// File: api/host_template_api/create.go
// Description: 主机模板创建

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// CreateRequest 创建主机模板的请求参数结构
type CreateRequest struct {
	Title    string                      `json:"title" binding:"required"` // 主机模板名称
	PortList models.HostTemplatePortList `json:"portList" binding:"dive"`  // 端口列表，需包含关联的虚拟服务ID和端口信息
}

// CreateView 创建主机模板的API入口函数
// 流程：1. 绑定并验证请求参数 2. 校验模板名称唯一性 3. 校验端口列表中端口是否重复 4. 校验关联的虚拟服务是否存在 5. 入库并返回结果
func (HostTemplateApi) CreateView(c *gin.Context) {
	// 绑定并验证请求参数（从上下文获取CreateRequest结构体）
	cr := middleware.GetBind[CreateRequest](c)

	// 校验主机模板名称是否已存在（名称不可重复）
	var model models.HostTemplateModel
	err := global.DB.Take(&model, "title = ? ", cr.Title).Error
	if err == nil {
		res.FailWithMsg("主机模板名称不能重复", c)
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

	// 构建主机模板模型并入库
	model = models.HostTemplateModel{
		Title:    cr.Title,
		PortList: cr.PortList,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithMsg("主机模板创建失败", c)
		return
	}

	// 返回创建成功的主机模板ID
	res.OkWithData(model.ID, c)
}
