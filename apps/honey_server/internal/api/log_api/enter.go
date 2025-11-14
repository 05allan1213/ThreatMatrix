package log_api

// File: api/log_api/enter.go
// Description: 日志管理的接口层，提供日志列表查询与日志删除功能

import (
	"fmt"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/common_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// LogApi 日志接口结构体
type LogApi struct {
}

// LogListRequest 日志列表查询参数
type LogListRequest struct {
	models.PageInfo
	Type int8   `form:"type"` // 日志类型（1 登录日志）
	IP   string `form:"ip"`   // 按 IP 精确搜索
	Addr string `form:"addr"` // 按归属地精确搜索
}

// LogListView 查询日志列表，支持分页、模糊搜索 username、按 IP/Addr 筛选
func (LogApi) LogListView(c *gin.Context) {
	cr := middleware.GetBind[LogListRequest](c) // 绑定查询参数

	// 根据筛选条件与分页参数查询列表
	list, count, _ := common_service.QueryList(models.LogModel{
		Type: cr.Type,
		IP:   cr.IP,
		Addr: cr.Addr,
	}, common_service.QueryListRequest{
		Likes:    []string{"username"}, // 用户名模糊搜索
		PageInfo: cr.PageInfo,          // 分页参数
		Sort:     "created_at desc",    // 按创建时间倒序
	})

	res.OkWithList(list, count, c)
}

// RemoveView 删除指定日志，支持批量删除
func (LogApi) RemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDListRequest](c) // 获取要删除的 ID 列表
	log := middleware.GetLog(c)                       // 当前操作人的日志信息

	// 执行删除操作
	successCount, err := common_service.Remove(models.LogModel{}, common_service.RemoveRequest{
		IDList:   cr.IdList,
		Log:      log,
		Msg:      "日志",
		Unscoped: true, // true 表示硬删除
	})

	if err != nil {
		msg := fmt.Sprintf("删除用户失败 %s", err)
		res.FailWithMsg(msg, c)
		return
	}

	msg := fmt.Sprintf("删除成功 共%d个，成功%d个", len(cr.IdList), successCount)
	res.OkWithMsg(msg, c)
}
