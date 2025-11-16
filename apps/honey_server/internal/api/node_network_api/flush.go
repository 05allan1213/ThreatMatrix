package node_network_api

// File: api/node_network_api/flush.go
// Description: 提供节点网卡信息刷新的API接口

import (
	"fmt"
	"honey_server/internal/rpc/node_rpc"
	"honey_server/internal/service/grpc_service"
	"honey_server/internal/utils/res"

	"github.com/gin-gonic/gin"
)

// 处理节点网卡信息刷新的API方法
func (NodeNetworkApi) FlushView(c *gin.Context) {
	// 向grpc命令请求通道发送网卡刷新指令
	grpc_service.CmdRequestChan <- &node_rpc.CmdRequest{
		CmdType: node_rpc.CmdType_cmdNetworkFlushType,
		TaskID:  "xx",
		NetworkFlushInMessage: &node_rpc.NetworkFlushInMessage{
			FilterNetworkName: []string{"hy-"},
		},
	}

	// 从grpc命令响应通道接收刷新结果
	response := <-grpc_service.CmdResponseChan
	fmt.Println("网卡刷新数据", response)

	// 返回成功响应，包含网卡刷新结果数据
	res.OkWithData(response.NetworkFlushOutMessage, c)
}
