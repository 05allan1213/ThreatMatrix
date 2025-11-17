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
	"time"

	"github.com/gin-gonic/gin"
)

// 处理节点网卡信息刷新的API方法
func (NodeNetworkApi) FlushView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	var model models.NodeModel
	err := global.DB.Take(&model, cr.Id).Error
	if err != nil {
		res.FailWithMsg("节点不存在", c)
		return
	}

	if model.Status != 1 {
		res.FailWithMsg("节点未运行", c)
		return
	}

	// 判断节点在不在
	_, ok := grpc_service.StreamMap[model.Uid]
	if !ok {
		res.FailWithMsg("节点离线中", c)
		return
	}
	// 向grpc命令请求通道发送网卡刷新指令
	grpc_service.CmdRequestChan <- &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  "xx",
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"},
		},
	}

	// 从grpc命令响应通道接收刷新结果，增加超时机制
	select {
	case response := <-grpc_service.CmdResponseChan:
		fmt.Println("网卡刷新数据", response)
		// 返回成功响应，包含网卡刷新结果数据
		res.OkWithData(response.NetworkFlushOutMessage, c)
	case <-time.After(30 * time.Second):
		// 超时处理
		fmt.Println("网卡刷新指令执行超时")
		res.FailWithMsg("节点响应超时", c)
	}
}
