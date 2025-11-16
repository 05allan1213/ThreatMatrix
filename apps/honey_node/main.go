package main

// File: main.go
// Description: gRPC客户端主程序

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"honey_node/internal/core"
	"honey_node/internal/global"
	"honey_node/internal/rpc/node_rpc"
	"honey_node/internal/utils/info"
	"honey_node/internal/utils/ip"
	"io/ioutil"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// 读取配置文件到全局配置变量
	global.Config = core.ReadConfig()
	// 设置默认日志配置
	core.SetLogDefault()
	// 获取日志实例
	global.Log = core.GetLogger()

	// 从配置中获取管理节点的gRPC服务地址
	addr := global.Config.System.GrpcManageAddr

	// 加载客户端证书和私钥
	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	if err != nil {
		logrus.Fatalf("failed to load client key pair: %v", err)
	}

	// 加载 CA 证书
	caCert, err := ioutil.ReadFile("cert/ca.crt")
	if err != nil {
		logrus.Fatalf("failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 创建 TLS 配置
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// 创建 credentials
	creds := credentials.NewTLS(config)

	// 创建 gRPC 连接
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		// 连接失败时打印错误并退出程序
		logrus.Fatalf("%s", fmt.Sprintf("grpc connect addr [%s] 连接失败 %s", addr, err))
	}
	// 延迟关闭连接，确保程序退出时释放资源
	defer conn.Close()

	// 初始化节点服务的gRPC客户端实例
	client := node_rpc.NewNodeServiceClient(conn)

	// 获取节点的IP地址和MAC地址
	_ip, mac, err := ip.GetNetworkInfo(global.Config.System.Network)
	if err != nil {
		logrus.Fatalln(err)
	}

	// 如果节点的UID为空，则生成一个新的UID并保存到配置文件中
	if global.Config.System.Uid == "" {
		global.Config.System.Uid = uuid.New().String()
		core.SetConfig()
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatalln(err)
	}

	// 发送节点注册请求到gRPC服务器
	_, err = client.Register(context.Background(), &node_rpc.RegisterRequest{
		Ip:      _ip,
		Mac:     mac,
		NodeUid: global.Config.System.Uid,
		Version: global.Version,
		Commit:  global.Commit,
		SystemInfo: &node_rpc.SystemInfoMessage{
			HostName: hostname,
		},
	})
	if err != nil {
		logrus.Fatalf("节点注册失败 %s", err)
		return
	}

	nodePath, _ := os.Getwd()
	fmt.Println(nodePath)
	resourceInfo, err := info.GetResourceInfo(nodePath)
	if err != nil {
		logrus.Fatalf("节点资源信息获取失败 %s", err)
		return
	}

	_, err = client.NodeResource(context.Background(), &node_rpc.NodeResourceRequest{
		NodeUid: global.Config.System.Uid,
		ResourceInfo: &node_rpc.ResourceMessage{
			CpuCount:              resourceInfo.CpuCount,
			CpuUseRate:            resourceInfo.CpuUseRate,
			MemTotal:              resourceInfo.MemTotal,
			MemUseRate:            resourceInfo.MemUseRate,
			DiskTotal:             resourceInfo.DiskTotal,
			DiskUseRate:           resourceInfo.DiskUseRate,
			NodePath:              resourceInfo.NodePath,
			NodeResourceOccupancy: resourceInfo.NodeResourceOccupancy,
		},
	})
	if err != nil {
		logrus.Fatalf("节点资源信息上报失败 %s", err)
		return
	}
}
