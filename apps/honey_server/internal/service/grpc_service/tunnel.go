package grpc_service

// File: service/grpc_service/tunnel.go
// Description: 实现双向流Tunnel接口，建立服务端与目标地址的TCP连接并转发数据，支持双向数据透传（隧道功能）

import (
	"fmt"
	"honey_server/internal/rpc/node_rpc"
	"io"
	"log"
	"net"
)

// Tunnel 实现gRPC双向流RPC接口，提供TCP隧道转发功能
// 客户端通过流发送目标地址及数据，服务端建立TCP连接并实现双向数据透传
func (s *NodeService) Tunnel(stream node_rpc.NodeService_TunnelServer) error {
	// 接收客户端的初始请求：获取需要连接的目标TCP地址（如"127.0.0.1:8080"）
	req, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("接收初始请求失败: %v", err)
	}

	// 创建TCP拨号器，使用stream的上下文（支持超时/取消）连接目标地址
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(stream.Context(), "tcp", req.Address)
	if err != nil {
		return fmt.Errorf("连接目标地址失败: %v", err)
	}
	defer conn.Close() // 函数退出时关闭TCP连接，释放资源

	// 启动goroutine处理"客户端→服务端→目标地址"的数据流向
	// 从gRPC流接收客户端数据，转发到目标TCP连接
	go func() {
		for {
			req, err := stream.Recv()
			if err == io.EOF { // 客户端关闭流则退出
				return
			}
			if err != nil {
				log.Printf("接收客户端数据失败: %v", err)
				return
			}

			// 将客户端发送的Chunk数据写入目标TCP连接
			_, err = conn.Write(req.Chunk)
			if err != nil {
				log.Printf("写入目标连接失败: %v", err)
				return
			}
		}
	}()

	// 处理"目标地址→服务端→客户端"的数据流向
	// 从目标TCP连接读取数据，通过gRPC流发送回客户端
	buffer := make([]byte, 4096) // 4KB缓冲区，平衡性能与内存占用
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("目标连接已关闭")
			} else {
				log.Printf("从目标连接读取失败: %v", err)
			}
			return nil // 目标连接关闭/出错时退出
		}

		// 将读取到的数据通过gRPC流发送给客户端
		err = stream.Send(&node_rpc.TunnelData{
			Chunk:   buffer[:n], // 仅发送实际读取到的字节（避免空数据）
			Address: req.Address,
		})
		if err != nil {
			log.Printf("发送数据到客户端失败: %v", err)
			return err
		}
	}
}
