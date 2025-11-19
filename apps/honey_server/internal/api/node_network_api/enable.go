package node_network_api

// File: api/node_network_api/enable.go
// Description: 节点网卡启用API

import (
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/utils/ip"
	"honey_server/internal/utils/res"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// mutex 用于并发控制的互斥锁，防止多请求同时操作同一网卡导致数据不一致
var mutex sync.Mutex

// EnableView 节点网卡启用接口处理函数
func (NodeNetworkApi) EnableView(c *gin.Context) {
	// 从请求中绑定并获取ID参数（models.IDRequest结构体）
	cr := middleware.GetBind[models.IDRequest](c)

	var model models.NodeNetworkModel
	// 根据ID查询网卡信息，并预加载关联的节点模型（NodeModel）
	err := global.DB.Preload("NodeModel").Take(&model, cr.Id).Error
	if err != nil {
		res.FailWithMsg("网卡不存在", c)
		return
	}

	// 加互斥锁，确保同一时间只有一个请求处理该网卡的启用操作
	mutex.Lock()
	defer mutex.Unlock()

	// 检查网卡状态，若已启用则直接返回错误
	if model.Status == 1 {
		res.FailWithMsg("网卡已启用，请勿重复启用", c)
		return
	}

	// 使用数据库事务执行原子操作，保证数据一致性
	err = global.DB.Transaction(func(tx *gorm.DB) error {
		// 解析CIDR格式获取可用IP范围（基于网卡的IP和掩码）
		ipRange, err1 := ip.ParseCIDRGetUseIPRange(fmt.Sprintf("%s/%d", model.IP, model.Mask))
		if err1 != nil {
			return err1 // 解析失败则回滚事务
		}

		// 构建网络模型数据，准备插入网络表
		var net = models.NetModel{
			NodeID:             model.NodeID,
			Title:              fmt.Sprintf("%s_%s_网络", model.NodeModel.Title, model.Network), // 拼接网络名称
			Network:            model.Network,
			IP:                 model.IP,
			Mask:               model.Mask,
			Gateway:            model.Gateway,
			CanUseHoneyIPRange: ipRange, // 可用诱捕IP范围
		}

		// 插入网络记录到数据库
		err = tx.Create(&net).Error
		if err != nil {
			return err // 插入失败则回滚事务
		}

		// 更新网卡状态为已启用（状态码1）
		err = tx.Model(&model).Update("status", 1).Error
		return err // 返回最后一步操作结果，决定事务提交或回滚
	})

	// 事务执行失败处理
	if err != nil {
		logrus.Errorf("网卡启用失败 %s", err) // 记录错误日志
		res.FailWithMsg("网卡启用失败", c)
		return
	}

	// 启用成功返回响应
	res.OkWithMsg("网卡启用成功", c)
}
