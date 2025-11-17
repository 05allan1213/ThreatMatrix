package node_network_api

// File: api/node_network_api/flush.go
// Description: 提供节点网卡信息刷新的API接口

import (
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// FlushView 处理指定节点的网卡信息刷新请求
func (NodeNetworkApi) FlushView(c *gin.Context) {
	// 获取并绑定请求参数（节点ID）
	cr := middleware.GetBind[models.IDRequest](c)

	// 根据ID查询节点信息，验证节点是否存在
	var model models.NodeModel
	err := global.DB.Take(&model, cr.Id).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}

	// 验证节点是否处于运行状态（状态1表示运行中）
	if model.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 检查节点是否在线（是否在grpc命令映射表中）
	_, ok := grpc_service.NodeCommandMap[model.Uid]
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}

	// 向目标节点的命令请求通道发送网卡刷新指令
	// 指令类型为网卡刷新，任务ID暂定为"xx"，过滤"hy-"前缀的网卡
	grpc_service.NodeCommandMap[model.Uid].ReqChan <- &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  "xx",
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"},
		},
	}

	// 从节点的命令响应通道接收刷新结果
	response := <-grpc_service.NodeCommandMap[model.Uid].ResChan
	fmt.Println("网卡刷新数据", response)

	// 返回成功响应，包含刷新后的网卡信息
	res.OkWithData(response.NetworkFlushOutMessage, c)
}
