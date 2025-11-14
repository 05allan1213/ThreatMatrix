package common_service

// File: service/common_service/query_list.go
// Description: 通用查询服务，支持精确匹配、模糊搜索、预加载、排序与分页。

import (
	"fmt"
	"honey_server/internal/core"
	"honey_server/internal/models"

	"gorm.io/gorm"
)

// QueryListRequest 通用查询参数结构体
type QueryListRequest struct {
	Debug    bool            // 是否开启 GORM Debug 模式
	Likes    []string        // 模糊查询字段列表
	Where    *gorm.DB        // 自定义查询条件
	Preload  []string        // 预加载字段
	Sort     string          // 排序规则
	PageInfo models.PageInfo // 分页信息
}

// QueryList 通用列表查询，支持预加载、模糊搜索、分页与排序
func QueryList[T any](model T, req QueryListRequest) (list []T, count int64, err error) {
	db := core.GetDB()
	if req.Debug {
		db = db.Debug() // 开启调试模式，打印 SQL
	}

	// 预加载字段
	for _, s := range req.Preload {
		db = db.Preload(s)
	}

	// 针对字段的精确匹配
	db = db.Where(model)

	// 高级查询条件（自定义 where）
	if req.Where != nil {
		db = db.Where(req.Where)
	}

	// 模糊匹配逻辑
	if req.PageInfo.Key != "" {
		like := core.GetDB().Where("")
		for _, column := range req.Likes {
			like.Or(fmt.Sprintf("%s like ?", column), fmt.Sprintf("%%%s%%", req.PageInfo.Key))
		}
		db = db.Where(like)
	}

	// 分页处理
	if req.PageInfo.Limit <= 0 {
		req.PageInfo.Limit = 10
	}
	if req.PageInfo.Page <= 0 {
		req.PageInfo.Page = 1
	}
	offset := (req.PageInfo.Page - 1) * req.PageInfo.Limit

	// 查询数据
	if err = db.Offset(offset).Limit(req.PageInfo.Limit).Order(req.Sort).Find(&list).Error; err != nil {
		return
	}

	// 查询总记录数
	err = db.Count(&count).Error
	return
}
