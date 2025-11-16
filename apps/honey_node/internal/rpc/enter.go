package rpc

// File: rpc/enter.go
// Description: 提供grpc客户端连接的创建功能，支持基于TLS证书的安全连接

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// 创建并返回一个grpc客户端连接
func GetConn(addr string) (conn *grpc.ClientConn) {
	// 加载客户端证书和私钥（用于TLS双向认证）
	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	if err != nil {
		logrus.Fatalf("加载客户端证书和私钥失败: %v", err)
	}

	// 加载CA根证书（用于验证服务端证书的合法性）
	caCert, err := ioutil.ReadFile("cert/ca.crt")
	if err != nil {
		logrus.Fatalf("读取CA根证书失败: %v", err)
	}
	// 创建CA证书池并添加CA根证书
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 创建TLS配置，指定客户端证书和信任的根证书池
	config := &tls.Config{
		Certificates: []tls.Certificate{cert}, // 客户端证书
		RootCAs:      caCertPool,              // 信任的根证书池
	}

	// 基于TLS配置创建grpc所需的凭证
	creds := credentials.NewTLS(config)

	// 使用指定的地址和凭证创建grpc客户端连接
	conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		logrus.Fatalf("%s", fmt.Sprintf("grpc连接地址 [%s] 失败: %s", addr, err))
	}
	return
}
