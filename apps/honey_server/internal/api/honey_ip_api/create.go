package honey_ip_api

// File: api/honey_ip_api/enter.go
// Description: 诱捕IP创建API

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// CreateRequest 诱捕IP创建请求结构体
type CreateRequest struct {
	NetID uint   `json:"netID" binding:"required"` // 网络ID
	IP    string `json:"ip" binding:"required"`    // 诱捕IP地址
}

// CreateView 诱捕IP创建接口处理函数
func (HoneyIPApi) CreateView(c *gin.Context) {
	// 从请求中绑定并获取创建参数（包含必填校验）
	cr := middleware.GetBind[CreateRequest](c)

	// 合法性校验1：判断网络是否存在（预加载关联的节点模型）
	var netModel models.NetModel
	err := global.DB.Preload("NodeModel").Take(&netModel, cr.NetID).Error
	if err != nil {
		res.FailWithMsg("网络不存在", c)
		return
	}

	// 合法性校验2：判断IP是否在网络的可用IP范围内
	ipRange, err := netModel.IpRange()
	if !utils.InList(ipRange, cr.IP) {
		res.FailWithMsg("当前ip不存在可部署ip列表里面", c)
		return
	}

	// 合法性校验3：判断IP是否已被主机占用
	var hostModel models.HostModel
	err = global.DB.Take(&hostModel, "net_id = ? and ip = ?", cr.NetID, cr.IP).Error
	if err == nil {
		res.FailWithMsg("当前ip是主机ip", c)
		return
	}

	// 合法性校验4：判断IP是否已被其他诱捕IP占用
	var honeyIPModel models.HoneyIpModel
	err = global.DB.Take(&honeyIPModel, "net_id = ? and ip = ?", cr.NetID, cr.IP).Error
	if err == nil {
		res.FailWithMsg("当前ip已使用", c)
		return
	}

	// 合法性校验5：判断节点状态是否为运行中
	if netModel.NodeModel.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 合法性校验6：通过gRPC检查节点是否在线
	_, ok := grpc_service.GetNodeCommand(netModel.NodeModel.Uid)
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 构建诱捕IP模型并入库
	var model = models.HoneyIpModel{
		NodeID: netModel.NodeID, // 关联节点ID
		NetID:  netModel.ID,     // 关联网络ID
		IP:     cr.IP,           // 诱捕IP地址
		Status: 1,               // 状态设为启用
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithMsg("创建诱捕ip失败", c)
		return
	}

	// 创建成功，返回诱捕IP记录ID
	res.OkWithData(model.ID, c)
}
