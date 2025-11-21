package honey_port_api

// File: api/honey_port_api/update.go
// Description: 诱捕转发端口更新API

import (
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/service/mq_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// UpdateRequest 诱捕端口更新请求结构体
type UpdateRequest struct {
	HoneyIPID uint       `json:"honeyIpID" binding:"required"`     // 关联的诱捕IP ID
	PortList  []PortType `json:"portList" binding:"dive,required"` // 端口配置列表
}

// PortType 端口配置项结构体
type PortType struct {
	Port      int  `json:"port" binding:"required,min=1,max=65535"` // 端口号
	ServiceID uint `json:"serviceID" binding:"required"`            // 关联的服务ID
}

// UpdateView 诱捕端口更新接口处理函数
func (HoneyPortApi) UpdateView(c *gin.Context) {
	// 从请求中绑定并获取更新参数（包含必填校验）
	cr := middleware.GetBind[UpdateRequest](c)

	// 校验关联的诱捕IP是否存在
	var honeyIPModel models.HoneyIpModel
	err := global.DB.Preload("NodeModel").Take(&honeyIPModel, cr.HoneyIPID).Error
	if err != nil {
		res.FailWithMsg("不存在的诱捕ip", c)
		return
	}

	nodeModel := honeyIPModel.NodeModel
	// 判断节点是否在线
	if nodeModel.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 获取节点的gRPC服务
	_, ok := grpc_service.GetNodeCommand(nodeModel.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 查询该诱捕IP下已配置的端口信息
	var honeyPortList []models.HoneyPortModel
	global.DB.Find(&honeyPortList, "honey_ip_id = ?", cr.HoneyIPID)

	// 端口配置合法性校验：
	// 1. 检查端口是否重复；2. 收集关联的服务ID用于后续有效性校验
	var portMap = map[int]struct{}{} // 用于检测端口重复的map
	var serviceIDList []uint         // 收集所有关联的服务ID
	for _, portType := range cr.PortList {
		serviceIDList = append(serviceIDList, portType.ServiceID)
		portMap[portType.Port] = struct{}{} // 记录端口号，用于检测重复
	}

	// 若端口map长度与请求端口列表长度不一致，说明存在重复端口
	if len(portMap) != len(cr.PortList) {
		res.FailWithMsg("端口重复", c)
		return
	}

	// 查询所有关联的服务信息，验证服务ID有效性并构建服务map
	var serviceList []models.ServiceModel
	global.DB.Find(&serviceList, "id in ?", serviceIDList)
	var serviceMap = map[uint]models.ServiceModel{}
	for _, model := range serviceList {
		serviceMap[model.ID] = model // 服务ID为key，便于快速查找
	}

	// 增量更新逻辑：计算需要新增和删除的端口
	// 1. 将现有端口转换为map（端口号为key），便于快速对比
	existingPorts := make(map[int]models.HoneyPortModel)
	for _, port := range honeyPortList {
		existingPorts[port.Port] = port
	}

	// 2. 计算需要新增的端口（请求中有但现有配置中没有的端口）
	var newPorts []models.HoneyPortModel
	for _, reqPort := range cr.PortList {
		// 校验服务ID是否有效
		service, ok := serviceMap[reqPort.ServiceID]
		if !ok {
			res.FailWithMsg(fmt.Sprintf("服务%d不存在", reqPort.ServiceID), c)
			return
		}

		// 若该端口未配置过，则加入新增列表
		if _, exists := existingPorts[reqPort.Port]; !exists {
			newPorts = append(newPorts, models.HoneyPortModel{
				HoneyIpID: cr.HoneyIPID,
				Port:      reqPort.Port,
				ServiceID: reqPort.ServiceID,
				DstIP:     service.IP,   // 从关联服务获取目标IP
				DstPort:   service.Port, // 从关联服务获取目标端口
				Status:    1,            // 假设1为启用状态
			})
		}
	}

	// 3. 计算需要删除的端口（现有配置中有但请求中没有的端口）
	var portsToDelete []models.HoneyPortModel
	for port, model := range existingPorts {
		found := false
		// 检查该端口是否存在于请求配置中
		for _, reqPort := range cr.PortList {
			if reqPort.Port == port {
				found = true
				break
			}
		}
		// 若请求中无该端口，则加入删除列表
		if !found {
			portsToDelete = append(portsToDelete, model)
		}
	}

	// 使用数据库事务执行端口的删除和新增操作，保证数据一致性
	tx := global.DB.Begin()
	if tx.Error != nil {
		res.FailWithMsg("更新端口信息失败", c)
		return
	}

	// 执行端口删除操作
	for _, port := range portsToDelete {
		if err := tx.Delete(&port).Error; err != nil {
			tx.Rollback() // 出错则回滚事务
			res.FailWithMsg("更新端口信息失败", c)
			return
		}
	}

	// 执行新端口添加操作
	for _, port := range newPorts {
		if err := tx.Create(&port).Error; err != nil {
			tx.Rollback() // 出错则回滚事务
			res.FailWithMsg("更新端口信息失败", c)
			return
		}
	}

	// 提交事务（所有操作成功才提交）
	if err := tx.Commit().Error; err != nil {
		res.FailWithMsg("更新端口信息失败", c)
		return
	}

	// 拼接更新结果信息并返回
	msg := fmt.Sprintf("新增端口%d个，删除端口%d个", len(newPorts), len(portsToDelete))

	// 查询当前诱捕IP下的所有端口信息，准备发送到消息队列
	var portList []models.HoneyPortModel
	global.DB.Find(&portList, "honey_ip_id = ?", cr.HoneyIPID)

	// 构造绑定端口请求消息
	req := mq_service.BindPortRequest{
		IP:    honeyIPModel.IP,
		LogID: "",
	}

	// 遍历端口列表，组装端口信息到请求中
	for _, model := range portList {
		req.PortList = append(req.PortList, mq_service.PortInfo{
			IP:       honeyIPModel.IP,
			Port:     model.Port,
			DestIP:   model.DstIP,
			DestPort: model.DstPort,
		})
	}

	// 发送绑定端口消息到对应的节点
	mq_service.SendBindPortMsg(nodeModel.Uid, req)

	res.OkWithMsg(msg, c)
}
