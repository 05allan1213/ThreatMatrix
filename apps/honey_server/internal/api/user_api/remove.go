package user_api

// File: api/user_api/remove.go
// Description: 用户批量删除接口，调用通用删除服务并返回处理结果。

import (
	"fmt"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// UserRemoveRequest 用户批量删除请求体
type UserRemoveRequest struct {
	IDList []uint `json:"idList"` // 待删除的用户 ID 列表
}

// UserRemoveView 用户删除接口
func (UserApi) UserRemoveView(c *gin.Context) {
	cr := middleware.GetBind[UserRemoveRequest](c) // 绑定参数
	log := middleware.GetLog(c)                    // 获取日志记录器

	// 调用通用删除服务
	successCount, err := common_service.Remove(
		models.UserModel{},
		common_service.RemoveRequest{
			IDList: cr.IDList, // 删除的 ID 列表
			Log:    log,       // 日志
			Msg:    "用户",      // 删除对象名称（用于日志）
		},
	)

	// 删除失败
	if err != nil {
		msg := fmt.Sprintf("删除用户失败 %s", err)
		res.FailWithMsg(msg, c)
		return
	}

	// 删除成功
	msg := fmt.Sprintf("删除成功 共%d个，成功%d个", len(cr.IDList), successCount)
	res.OkWithMsg(msg, c)
}
