package matrix_template_api

// File: api/matrix_template_api/create.go
// Description: 矩阵模板创建

import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// CreateRequest 创建矩阵模板的请求参数结构
type CreateRequest struct {
	Title            string                    `json:"title" binding:"required"`                          // 矩阵模板名称
	HostTemplateList []models.HostTemplateInfo `json:"hostTemplateList" binding:"required,dive,required"` // 关联的主机模板列表
}

// CreateView 创建矩阵模板的API入口函数
//
// 1. 绑定并验证请求参数
// 2. 校验主机模板列表不为空
// 3. 校验模板名称唯一性
// 4. 收集关联的主机模板ID并验证其存在性
// 5. 入库并返回结果
func (MatrixTemplateApi) CreateView(c *gin.Context) {
	// 绑定并验证请求参数（从上下文获取CreateRequest结构体）
	cr := middleware.GetBind[CreateRequest](c)

	// 校验主机模板列表不能为空（矩阵模板至少需关联一个主机模板）
	if len(cr.HostTemplateList) == 0 {
		res.FailWithMsg("矩阵模板需要关联至少一个主机模板", c)
		return
	}

	// 校验矩阵模板名称是否已存在（名称不可重复）
	var model models.MatrixTemplateModel
	err := global.DB.Take(&model, "title = ? ", cr.Title).Error
	if err == nil {
		res.FailWithMsg("矩阵模板名称不能重复", c)
		return
	}

	// 收集主机模板列表中所有关联的主机模板ID，用于后续存在性校验
	var hostTemplateIDList []uint
	for _, h := range cr.HostTemplateList {
		hostTemplateIDList = append(hostTemplateIDList, h.HostTemplateID)
	}

	// 查询所有关联的主机模板，验证是否存在
	var hostTemps []models.HostTemplateModel
	global.DB.Find(&hostTemps, "id in ?", hostTemplateIDList)
	// 将查询到的主机模板按ID映射，便于快速查找
	var hostTempMap = map[uint]models.HostTemplateModel{}
	for _, m := range hostTemps {
		hostTempMap[m.ID] = m
	}

	// 校验主机模板列表中关联的主机模板是否均存在
	for _, h := range cr.HostTemplateList {
		_, ok := hostTempMap[h.HostTemplateID]
		if !ok {
			msg := fmt.Sprintf("主机模板 %d 不存在", h.HostTemplateID)
			res.FailWithMsg(msg, c)
			return
		}
	}

	// 构建矩阵模板模型并入库
	model = models.MatrixTemplateModel{
		Title:            cr.Title,
		HostTemplateList: cr.HostTemplateList,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithMsg("矩阵模板创建失败", c)
		return
	}

	// 返回创建成功的矩阵模板ID
	res.OkWithData(model.ID, c)
}
