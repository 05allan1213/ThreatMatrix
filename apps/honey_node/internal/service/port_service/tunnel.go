package port_service

// File: service/port_service/tunnel.go
// Description: 实现本地TCP端口监听与gRPC隧道的桥接，将本地端口请求转发到服务端指定的目标地址，支撑诱捕端口的数据透传功能

import (
	"context"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var tunnelStore = sync.Map{}

// Tunnel 启动本地TCP端口监听并建立gRPC隧道转发
// 实现本地端口到目标地址的TCP数据透传，支撑诱捕端口的代理功能
func Tunnel(localAddr, targetAddr string) (err error) {
	// 创建本地TCP监听：绑定指定的本地地址和端口
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		logrus.Errorf("创建本地监听失败: %v", err)
		return
	}

	logrus.Infof("本地监听启动，地址: %s", localAddr)
	logrus.Infof("目标地址: %s", targetAddr)

	tunnelStore.Store(localAddr, listener)

	// 循环接受客户端连接：持续监听端口，处理新的TCP连接请求
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				break
			}
			logrus.Errorf("接受客户端连接失败: %v", err)
			break
		}

		// 异步处理单个连接的转发逻辑：每个连接使用独立goroutine，支持高并发
		go handleConnection(global.GrpcClient, clientConn, targetAddr)
	}
	return nil
}

// handleConnection 处理单个TCP连接的gRPC隧道转发
// 通过gRPC双向流将本地连接数据与服务端目标地址数据进行透传
func handleConnection(client node_rpc.NodeServiceClient, localConn net.Conn, targetAddr string) {
	defer localConn.Close() // 连接处理完成后关闭本地TCP连接，释放资源

	// 创建带取消功能的上下文：用于控制gRPC流的生命周期，支持连接中断时的资源清理
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 函数退出时取消上下文，终止gRPC流

	// 创建gRPC Tunnel双向流：建立与服务端的隧道连接，用于数据透传
	stream, err := client.Tunnel(ctx)
	if err != nil {
		logrus.Infof("创建隧道失败: %v", err)
		return
	}

	// 发送初始消息：向服务端传递目标转发地址，触发服务端连接目标地址
	if err := stream.Send(&node_rpc.TunnelData{
		Chunk:   []byte{}, // 初始消息无业务数据，仅传递地址信息
		Address: targetAddr,
	}); err != nil {
		logrus.Errorf("发送初始请求失败: %v", err)
		return
	}

	// 启动goroutine处理"服务端→gRPC流→本地TCP连接"的数据流向
	// 从gRPC流读取服务端转发的目标地址数据，写入本地客户端连接
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF { // 服务端关闭流时退出循环
				break
			}
			if err != nil {
				logrus.Errorf("接收gRPC服务器数据失败: %v", err)
				break
			}

			// 将服务端转发的数据写入本地TCP连接，完成"服务端→客户端"的数据传输
			_, err = localConn.Write(resp.Chunk)
			if err != nil {
				logrus.Errorf("写入本地连接失败: %v", err)
				break
			}
		}
		cancel() // 数据流向中断时取消上下文，终止另一方向的处理
	}()

	// 处理"本地TCP连接→gRPC流→服务端"的数据流向
	// 从本地客户端连接读取数据，通过gRPC流发送到服务端目标地址
	buffer := make([]byte, 4096) // 4KB缓冲区：平衡IO性能与内存占用，适配多数网络数据包大小
	for {
		n, err := localConn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				logrus.Infof("本地连接已关闭")
			} else {
				logrus.Errorf("从本地连接读取失败: %v", err)
			}
			break // 本地连接关闭或读取失败时退出循环
		}

		// 将本地连接读取的数据通过gRPC流发送到服务端，完成"客户端→服务端"的数据传输
		err = stream.Send(&node_rpc.TunnelData{
			Chunk:   buffer[:n], // 仅发送实际读取的有效字节，避免空数据传输
			Address: targetAddr,
		})
		if err != nil {
			logrus.Errorf("发送数据到gRPC服务器失败: %v", err)
			break
		}
	}

	// 主动关闭gRPC流的发送端：告知服务端本地数据发送完成，触发服务端的流关闭逻辑
	stream.CloseSend()
}

// CloseIpTunnel 关闭指定IP的TCP监听
func CloseIpTunnel(ip string) {
	tunnelStore.Range(func(key, value any) bool {
		localAddr := key.(string)
		if strings.HasPrefix(localAddr, ip) {
			logrus.Infof("清除%s上的全部服务", ip)
			listener := value.(net.Listener)
			listener.Close()
		}
		return true
	})
}
