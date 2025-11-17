package node_network_api

// File: api/node_network_api/flush.go
// Description: 提供节点网卡信息刷新与数据库同步的API

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// FlushView 刷新节点网卡信息并同步到数据库
func (NodeNetworkApi) FlushView(c *gin.Context) {
	// 获取并绑定请求参数（节点ID），通过中间件进行参数校验
	cr := middleware.GetBind[models.IDRequest](c)

	// 根据ID查询节点信息，验证节点是否存在
	var model models.NodeModel
	if err := global.DB.Take(&model, cr.Id).Error; err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}

	// 验证节点是否处于运行状态（状态1表示运行中）
	if model.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 通过封装的方法获取节点对应的命令交互实例，验证节点是否在线
	cmd, ok := grpc_service.GetNodeCommand(model.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 构建网卡刷新命令请求，使用当前时间戳的纳秒数作为唯一任务ID
	req := &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,           // 命令类型：网卡刷新
		TaskID:  fmt.Sprintf("flush-%d", time.Now().UnixNano()), // 任务ID，用于关联命令与响应
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"}, // 过滤"hy-"前缀的网卡（如诱捕相关网卡）
		},
	}

	// 创建带30秒超时的上下文，控制命令发送和响应接收的超时时间
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel() // 确保上下文最终被取消，释放资源

	// 向节点的命令请求通道发送刷新命令，带超时控制
	select {
	case cmd.ReqChan <- req:
		logrus.Debugf("已向节点[%s]发送网卡刷新请求", model.Uid)
	case <-ctx.Done():
		res.FailWithMsg("发送命令超时", c)
		return
	}

	// 存储从节点获取的最新网卡列表
	var networkInfoList []*node_rpc.NetworkInfoMessage

	// 等待节点返回的刷新结果，带超时控制
	select {
	case response := <-cmd.ResChan:
		logrus.Debugf("已接收节点[%s]的网卡刷新响应", model.Uid)
		networkInfoList = response.NetworkFlushOutMessage.NetworkList // 提取最新网卡列表
	case <-ctx.Done():
		res.FailWithMsg("获取响应超时", c)
		return
	}

	// 查询数据库中当前节点的网卡记录（用于后续对比）
	var networkList []models.NodeNetworkModel
	global.DB.Find(&networkList, "node_id = ?", model.ID)

	// 构建网卡名称到索引的映射（优化查找效率，避免嵌套循环）
	// 数据库中已存在的网卡映射：key=网卡名称，value=在networkList中的索引
	networkMap := make(map[string]int)
	for i, network := range networkList {
		networkMap[network.Network] = i
	}

	// 节点返回的新网卡列表映射：key=网卡名称，value=在networkInfoList中的索引
	newNetworkMap := make(map[string]int)
	for i, network := range networkInfoList {
		logrus.Infof("节点返回的网卡信息：名称=%s，网络地址=%s", network.Network, network.Net)
		newNetworkMap[network.Network] = i
	}

	// 计算新增的网卡（节点有但数据库没有的）
	var newNetworks []*node_rpc.NetworkInfoMessage
	for networkName := range newNetworkMap {
		if _, exists := networkMap[networkName]; !exists {
			// 从节点返回的列表中取出该网卡信息
			newNetworks = append(newNetworks, networkInfoList[newNetworkMap[networkName]])
		}
	}

	// 计算需要删除的网卡（数据库有但节点没有的）
	var deletedNetworks []models.NodeNetworkModel
	for networkName := range networkMap {
		if _, exists := newNetworkMap[networkName]; !exists {
			// 从数据库记录中取出该网卡信息
			deletedNetworks = append(deletedNetworks, networkList[networkMap[networkName]])
		}
	}

	// 计算需要更新的网卡（名称存在但IP或掩码有变化的）
	var updatedNetworks []models.NodeNetworkModel
	for networkName := range networkMap {
		if newIndex, exists := newNetworkMap[networkName]; exists {
			// 获取数据库中的记录和节点返回的新记录
			dbNetwork := networkList[networkMap[networkName]]
			newNetwork := networkInfoList[newIndex]

			// 检查IP或掩码是否发生变化
			if dbNetwork.IP != newNetwork.Ip || dbNetwork.Mask != int8(newNetwork.Mask) {
				// 更新数据库记录的IP和掩码
				dbNetwork.IP = newNetwork.Ip
				dbNetwork.Mask = int8(newNetwork.Mask)
				updatedNetworks = append(updatedNetworks, dbNetwork)
			}
		}
	}

	// 执行数据库操作：新增网卡记录
	for _, network := range newNetworks {
		newRecord := models.NodeNetworkModel{
			NodeID:  model.ID,           // 关联的节点ID
			Network: network.Network,    // 网卡名称
			IP:      network.Ip,         // 网卡IP地址
			Mask:    int8(network.Mask), // 子网掩码
			Status:  2,                  // 状态：2表示未启用
		}
		if err := global.DB.Create(&newRecord).Error; err != nil {
			logrus.Errorf("新增网卡记录失败: %v", err)
		}
	}

	// 执行数据库操作：删除无效网卡记录
	for _, network := range deletedNetworks {
		if err := global.DB.Delete(&network).Error; err != nil {
			logrus.Errorf("删除网卡记录失败: %v", err)
		}
	}

	// 执行数据库操作：更新变化的网卡记录
	for _, network := range updatedNetworks {
		if err := global.DB.Save(&network).Error; err != nil {
			logrus.Errorf("更新网卡记录失败: %v", err)
		}
	}

	// 记录同步结果日志
	logrus.Infof("网卡信息同步完成 - 新增: %d 个, 删除: %d 个, 更新: %d 个", len(newNetworks), len(deletedNetworks), len(updatedNetworks))

	// 返回成功响应
	res.OkWithMsg("网卡信息更新成功", c)
}
