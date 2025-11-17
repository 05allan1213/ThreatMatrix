package node_network_api

// File: api/node_network_api/enable.go
// Description: 节点网卡启用API

import (
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// EnableView 处理启用网卡的请求
func (n *NodeNetworkApi) EnableView(c *gin.Context) {
	// 从请求中获取绑定的ID参数（用于指定要启用的网卡）
	cr := middleware.GetBind[models.IDRequest](c)

	// 查询指定ID的网卡信息，并关联查询对应的节点信息
	var model models.NodeNetworkModel
	err := global.DB.Preload("NodeModel").Take(&model, cr.Id).Error
	if err != nil {
		// 若查询失败（网卡不存在），返回错误提示
		res.FailWithMsg("网卡不存在", c)
		return
	}

	// 加互斥锁防止并发操作冲突，确保启用过程的原子性
	n.mutex.Lock()
	defer n.mutex.Unlock()

	// 检查网卡当前状态，若已启用则返回错误提示（避免重复启用）
	if model.Status == 1 {
		res.FailWithMsg("网卡已启用，请勿重复启用", c)
		return
	}

	// 使用数据库事务确保操作的原子性，避免部分操作成功导致数据不一致
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 创建网络表记录，关联节点和网卡的网络信息
		title := fmt.Sprintf("%s_%s_网络", model.NodeModel.Title, model.Network)
		// 确保title不超过32个字符，避免数据库错误
		if len(title) > 32 {
			// 截取前面的部分，保证至少包含一部分网络名和"_网络"后缀
			title = title[:32]
		}
		var net = models.NetModel{
			NodeID:  model.NodeID,
			Title:   title,
			Network: model.Network,
			IP:      model.IP,
			Mask:    model.Mask,
			Gateway: model.Gateway,
		}
		err = tx.Create(&net).Error
		if err != nil {
			return err // 网络记录创建失败，触发事务回滚
		}

		// 创建主机表记录，关联节点和对应的网络信息
		var host = models.HostModel{
			NodeID: model.NodeID,
			NetID:  net.ID,
			IP:     net.IP,
			//Mac （预留字段，暂未赋值）
			//Manuf （预留字段，暂未赋值）
		}
		err = tx.Create(&host).Error
		if err != nil {
			return err // 主机记录创建失败，触发事务回滚
		}

		// 更新网卡状态为启用
		err = tx.Model(&model).Update("status", 1).Error
		return err // 返回状态更新结果，若失败则触发事务回滚
	})

	// 处理事务执行结果
	if err != nil {
		// 事务执行失败，记录错误日志并返回提示
		logrus.Errorf("网卡启用失败 %s", err)
		res.FailWithMsg("网卡启用失败", c)
		return
	}

	// 所有操作成功，返回启用成功提示
	res.OkWithMsg("网卡启用成功", c)
}
