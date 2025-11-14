package user_api

// File: api/user_api/list.go
// Description: 用户分页查询接口

import (
	"honey_server/middleware"
	"honey_server/models"
	"honey_server/service/common_service"
	"honey_server/utils/res"

	"github.com/gin-gonic/gin"
)

// UserListRequest 定义用户分页查询请求参数
type UserListRequest struct {
	models.PageInfo
	Username string `form:"username"` // 用户名查询参数
}

// UserListView 用户分页查询接口
func (UserApi) UserListView(c *gin.Context) {
	// 绑定请求参数并自动校验
	cr := middleware.GetBind[UserListRequest](c)

	// 查询用户列表，支持用户名模糊查询、分页和时间倒序排序
	list, count, _ := common_service.QueryList(
		models.UserModel{Username: cr.Username},
		common_service.Request{
			Likes:    []string{"username"}, // 模糊查询字段
			PageInfo: cr.PageInfo,          // 分页信息
			Sort:     "created_at desc",    // 排序规则
		},
	)

	// 返回带总数的列表响应
	res.OkWithList(list, count, c)
}
