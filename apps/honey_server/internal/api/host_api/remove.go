package host_api

// File: api/host_api/remove.go
// Description: 主机删除API
import (
	"fmt"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// RemoveView 主机批量删除接口处理函数
func (HostApi) RemoveView(c *gin.Context) {
	// 从请求中绑定并获取批量删除的ID列表参数（models.IDListRequest结构体）
	cr := middleware.GetBind[models.IDListRequest](c)
	// 从上下文获取日志实例，用于记录删除操作相关日志
	log := middleware.GetLog(c)

	// 调用通用删除服务执行主机批量删除
	successCount, err := common_service.Remove(models.HostModel{}, common_service.RemoveRequest{
		IDList: cr.IdList,
		Log:    log,
		Msg:    "主机",
	})

	// 删除过程中出现错误则返回失败信息
	if err != nil {
		msg := fmt.Sprintf("删除主机失败 %s", err)
		res.FailWithMsg(msg, c)
		return
	}

	// 拼接删除成功的结果信息（总数、成功数）并返回
	msg := fmt.Sprintf("删除成功 共%d个，成功%d个", len(cr.IdList), successCount)
	res.OkWithMsg(msg, c)
}
