package net_api

// File: api/net_api/scan.go
// Description: 网络扫描API

import (
	"context"
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/ip"
	"honey_server/internal/utils/res"
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var mutex sync.Mutex
var netProgressMap sync.Map

// ScanView 处理网络扫描任务的请求
// 发起扫描命令到节点，立即返回任务状态，异步处理扫描结果并更新数据库
func (NetApi) ScanView(c *gin.Context) {
	// 获取并绑定请求参数（包含要扫描的网络ID）
	cr := middleware.GetBind[models.IDRequest](c)

	// 查询指定ID的网络信息，并预加载关联的节点信息
	var model models.NetModel
	if err := global.DB.Preload("NodeModel").Take(&model, cr.Id).Error; err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}

	// 检查网络所属节点是否处于运行状态（状态1为运行）
	if model.NodeModel.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 计算可用IP范围（排除网络地址和广播地址）
	if model.CanUseHoneyIPRange == "" {
		// 计算起始ip
		_, ipNet, err := net.ParseCIDR(model.Subnet())
		if err != nil {
			res.FailWithMsg("无效的子网格式", c)
			return
		}

		// 计算可用IP范围（排除网络地址和广播地址）
		startIP := ip.IncrementIP(ipNet.IP)
		endIP := ip.DecrementIP(ip.BroadcastIP(ipNet))
		model.CanUseHoneyIPRange = ip.FormatIPRange(startIP, endIP)
	}

	// 通过节点唯一标识（Uid）获取节点命令通道（用于发送grpc命令和接收响应）
	cmd, ok := grpc_service.GetNodeCommand(model.NodeModel.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 过滤诱捕ip
	var filterIPList []string
	global.DB.Model(models.HoneyIpModel{}).Where("net_id = ?", cr.Id).Select("ip").Scan(&filterIPList)
	fmt.Println("过滤的ip列表", filterIPList)

	// 生成唯一任务ID（基于当前时间戳的纳秒级）
	taskID := fmt.Sprintf("netScan-%d", time.Now().UnixNano())
	// 构建网络扫描请求参数
	req := &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetScanType, // 命令类型：网络扫描
		TaskID:  taskID,                          // 任务唯一标识
		NetScanInMessage: &node_rpc.NetScanInMessage{
			Network:      model.Network,            // 目标网络接口
			IpRange:      model.CanUseHoneyIPRange, // 扫描的IP范围
			FilterIPList: filterIPList,             // 过滤IP列表
			NetID:        uint32(model.ID),         // 网络ID
		},
	}

	mutex.Lock()
	if model.ScanStatus == 2 {
		res.FailWithMsg("当前子网正在扫描中", c)
		mutex.Unlock()
		return
	}

	// 修改状态 为扫描中
	global.DB.Model(&model).Update("scan_status", 2)
	mutex.Unlock()

	// 发送扫描请求到节点的命令通道（非阻塞发送，避免通道繁忙时阻塞）
	select {
	case cmd.ReqChan <- req:
		logrus.Debugf("已向节点 %s 发送扫描请求，任务ID: %s", model.NodeModel.Uid, taskID)
	default:
		res.FailWithMsg("发送命令通道繁忙", c)
		return
	}

	// 立即返回给客户端，告知扫描任务已启动（不等待扫描完成）
	res.Ok(map[string]string{
		"task_id": taskID,
		"message": "扫描任务已启动，请稍后查询结果",
	}, "扫描任务已启动", c)

	// 异步处理扫描结果（不阻塞当前请求响应）
	go func(nodeUid string, netModel models.NetModel, cmdChan *grpc_service.Command, taskID string) {
		// 为异步处理创建独立上下文，设置5分钟超时（适应长时间扫描）
		ctxAsync, cancelAsync := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancelAsync()

		// 收集扫描过程中的有效结果
		var netScanMsg []*node_rpc.NetScanOutMessage
	label:
		for {
			select {
			case response := <-cmdChan.ResChan:
				// 过滤非当前任务的响应（确保只处理本任务结果）
				if response.TaskID != taskID {
					// 将非当前任务的响应重新放入通道，避免丢失
					select {
					case cmdChan.ResChan <- response:
					case <-ctxAsync.Done():
						break label
					}
					continue
				}

				logrus.Debugf("已接收节点 %s 的扫描响应，任务ID: %s", nodeUid, taskID)
				message := response.NetScanOutMessage

				// 若扫描过程中出现错误，终止循环
				if message.ErrMsg != "" {
					logrus.Errorf("节点 %s 扫描错误: %s", nodeUid, message.ErrMsg)
					break label
				}

				// 若扫描结束，终止循环
				if message.End {
					break label
				}

				// 收集包含有效IP的扫描结果
				if message.Ip != "" {
					netScanMsg = append(netScanMsg, message)
					netProgressMap.Store(uint(message.NetID), float64(message.Progress))
					fmt.Printf("网络扫描 %s %s %s %.2f\n", message.Ip, message.Mac, message.Manuf, message.Progress)
				}

			case <-ctxAsync.Done():
				// 异步处理超时
				logrus.Errorf("节点 %s 扫描超时，任务ID: %s", nodeUid, taskID)
				return
			}
		}

		// 处理扫描结果（若有有效结果且未超时）
		if len(netScanMsg) > 0 || ctxAsync.Err() == nil {
			processScanResult(netModel, netScanMsg)
		} else {
			logrus.Warnf("任务 %s 未接收到有效扫描结果", taskID)
		}

	}(model.NodeModel.Uid, model, cmd, taskID)
}

// processScanResult 处理扫描结果，对比数据库中已有的主机信息，更新新增、变更和删除的主机记录
func processScanResult(netModel models.NetModel, scanMsgs []*node_rpc.NetScanOutMessage) {
	// 在函数执行完毕后，将扫描状态更新为完成，并从进度映射中删除该网络
	defer func() {
		global.DB.Model(&netModel).Updates(map[string]any{
			"scan_progress": 100,
			"scan_status":   1,
		})
		netProgressMap.Delete(netModel.ID)
	}()

	// 查询当前网络下的所有主机信息
	var hostList []models.HostModel
	if err := global.DB.Find(&hostList, "net_id = ?", netModel.ID).Error; err != nil {
		logrus.Errorf("查询网络 %d 的主机列表失败: %v", netModel.ID, err)
		return
	}

	// 1. 将数据库中的主机列表转换为以IP为键的映射（便于快速查询）
	dbHostMap := make(map[string]models.HostModel)
	for _, host := range hostList {
		dbHostMap[host.IP] = host
	}

	// 2. 将扫描结果转换为以IP为键的映射（便于对比）
	scanResultMap := make(map[string]*node_rpc.NetScanOutMessage)
	for _, msg := range scanMsgs {
		if msg.Ip != "" {
			scanResultMap[msg.Ip] = msg
		}
	}

	// 3. 对比两个映射，确定新增、更新和删除的主机
	var newHosts []models.HostModel     // 新增的主机
	var deletedHostIDs []uint           // 需删除的主机ID
	var updatedHosts []models.HostModel // 需更新的主机

	// 处理新增和更新的主机
	for ip, scanMsg := range scanResultMap {
		if dbHost, exists := dbHostMap[ip]; exists {
			// 主机已存在，检查MAC或厂商信息是否变更
			if dbHost.Mac != scanMsg.Mac || dbHost.Manuf != scanMsg.Manuf {
				dbHost.Mac = scanMsg.Mac
				dbHost.Manuf = scanMsg.Manuf
				updatedHosts = append(updatedHosts, dbHost)
			}
			delete(dbHostMap, ip) // 移除已处理的主机，剩余的即为需删除的
		} else {
			// 主机不存在，新增记录
			newHosts = append(newHosts, models.HostModel{
				NodeID: netModel.NodeModel.ID,
				NetID:  netModel.ID,
				IP:     scanMsg.Ip,
				Mac:    scanMsg.Mac,
				Manuf:  scanMsg.Manuf,
			})
		}
	}

	// 收集需删除的主机ID（数据库中存在但扫描结果中不存在的主机）
	for _, dbHost := range dbHostMap {
		deletedHostIDs = append(deletedHostIDs, dbHost.ID)
	}

	// 打印扫描结果统计信息
	logrus.Infof("网络 %d 扫描结果：新增=%d，更新=%d，删除=%d",
		netModel.ID, len(newHosts), len(updatedHosts), len(deletedHostIDs))

	// 4. 使用事务批量更新数据库（确保操作原子性）
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 新增主机
		if len(newHosts) > 0 {
			if err := tx.Create(&newHosts).Error; err != nil {
				return fmt.Errorf("创建新主机失败: %w", err)
			}
		}

		// 更新主机信息
		if len(updatedHosts) > 0 {
			for _, host := range updatedHosts {
				if err := tx.Model(&models.HostModel{}).
					Where("id = ?", host.ID).
					Updates(map[string]interface{}{
						"mac":   host.Mac,
						"manuf": host.Manuf,
					}).Error; err != nil {
					return fmt.Errorf("更新主机信息失败: %w", err)
				}
			}
		}

		// 删除主机
		if len(deletedHostIDs) > 0 {
			if err := tx.Delete(&models.HostModel{}, deletedHostIDs).Error; err != nil {
				return fmt.Errorf("删除主机失败: %w", err)
			}
		}

		return nil
	})

	// 记录事务执行结果
	if err != nil {
		logrus.Errorf("更新网络 %d 的扫描结果失败: %v", netModel.ID, err)
	} else {
		logrus.Infof("成功更新网络 %d 的扫描结果", netModel.ID)
	}
}
