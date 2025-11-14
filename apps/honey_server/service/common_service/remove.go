package common_service

// File: service/common_service/remove.go
// Description: 通用删除服务，支持按条件、ID 列表删除，并记录详细日志。

import (
	"honey_server/global"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// RemoveRequest 通用删除请求结构
type RemoveRequest struct {
	Debug    bool          // 是否开启 Debug 模式
	Where    *gorm.DB      // 自定义 where 条件
	IDList   []uint        // ID 列表，用于批量删除
	Log      *logrus.Entry // 日志记录器
	Msg      string        // 删除对象名称，用于日志提示
	Unscoped bool          // 是否进行硬删除
}

// Remove 通用删除方法，支持 ID 列表、条件删除、日志记录
func Remove[T any](model T, req RemoveRequest) (successCount int64, err error) {
	db := global.DB
	deleteDB := global.DB

	// Debug 打印 SQL
	if req.Debug {
		db = db.Debug()
		deleteDB = deleteDB.Debug()
	}

	// 启用硬删除
	if req.Unscoped {
		req.Log.Infof("启用真删除")
		deleteDB = deleteDB.Unscoped()
	}

	// 添加自定义查询条件
	if req.Where != nil {
		db = db.Where(req.Where)
	}

	// 精确匹配传入的 model 字段
	db = db.Where(model)

	// 根据 ID 列表删除
	if len(req.IDList) > 0 {
		req.Log.Infof("删除 %s idList %v", req.Msg, req.IDList)
		db = db.Where("id in ?", req.IDList)
	}

	// 查询目标记录
	var list []T
	db.Find(&list)

	// 若没有查到记录
	if len(list) <= 0 {
		req.Log.Infof("没查到")
		return
	}

	// 执行删除
	result := deleteDB.Delete(&list)
	if result.Error != nil {
		req.Log.Errorf("删除失败 %s", result.Error)
		return
	}

	// 删除成功数量
	successCount = result.RowsAffected
	req.Log.Infof("删除 %s 成功, 成功%d个", req.Msg, successCount)
	return
}
