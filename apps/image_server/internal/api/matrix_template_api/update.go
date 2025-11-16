package matrix_template_api

// File: api/matrix_template_api/update.go
// Description: 矩阵模板更新接口
import (
	"fmt"
	"image_server/internal/global"
	"image_server/internal/middleware"
	"image_server/internal/models"
	"image_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// UpdateRequest 更新矩阵模板的请求参数结构
type UpdateRequest struct {
	ID               uint                      `json:"id" binding:"required"`                             // 矩阵模板ID
	Title            string                    `json:"title" binding:"required"`                          // 更新后的模板名称
	HostTemplateList []models.HostTemplateInfo `json:"hostTemplateList" binding:"required,dive,required"` // 更新后的关联主机模板列表
}

// UpdateView 更新矩阵模板的API入口函数
//
// 1. 绑定并验证请求参数
// 2. 检查待更新的模板是否存在
// 3. 校验主机模板列表不为空
// 4. 校验新名称是否与其他模板重复
// 5. 校验关联的主机模板是否存在
// 6. 执行更新并返回结果
func (MatrixTemplateApi) UpdateView(c *gin.Context) {
	// 绑定并验证请求参数
	cr := middleware.GetBind[UpdateRequest](c)

	// 检查待更新的矩阵模板是否存在
	var model models.MatrixTemplateModel
	err := global.DB.Take(&model, cr.ID).Error
	if err != nil {
		res.FailWithMsg("矩阵模板不存在", c)
		return
	}

	// 校验主机模板列表不能为空（矩阵模板至少需关联一个主机模板）
	if len(cr.HostTemplateList) == 0 {
		res.FailWithMsg("矩阵模板需要关联至少一个主机模板", c)
		return
	}

	// 校验新模板名称是否与其他模板重复（排除当前模板自身）
	var newModel models.MatrixTemplateModel
	err = global.DB.Take(&newModel, "title = ? and id <> ?", cr.Title, cr.ID).Error
	if err == nil {
		res.FailWithMsg("修改的矩阵模板名称不能重复", c)
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

	// 构建更新数据并执行数据库更新操作
	newModel = models.MatrixTemplateModel{
		Title:            cr.Title,
		HostTemplateList: cr.HostTemplateList,
	}
	err = global.DB.Model(&model).Updates(newModel).Error
	if err != nil {
		res.FailWithMsg("矩阵模板修改失败", c)
		return
	}

	// 返回更新成功消息
	res.OkWithMsg("矩阵模板修改成功", c)
}
