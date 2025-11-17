package ip

// File:utils/ip/enter.go
// Description: 判断ip是否是本地ip

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// 判断ip是否是本地ip
func HasLocalIPAddr(_ip string) bool {
	ip := net.ParseIP(_ip)
	// 判断是否是私有地址
	if ip.IsPrivate() {
		return true
	}
	// 判断是否是回环地址
	if ip.IsLoopback() {
		return true
	}
	return false
}

// 解析IP范围字符串，支持单个IP和IP段格式，返回IP字符串列表
func ParseIPRange(ipRange string) ([]string, error) {
	var result []string
	// 按逗号分割字符串，支持同时解析多个IP/IP段（如"192.168.1.1,192.168.1.5-8"）
	segments := strings.Split(ipRange, ",")

	// 遍历每个分割后的IP/IP段进行处理
	for _, segment := range segments {
		// 去除首尾空格
		segment = strings.TrimSpace(segment)
		if segment == "" {
			// 跳过空字符串（如连续逗号导致的空项）
			continue
		}

		// 检查是否包含连字符，判断是否为IP段（而非单个IP）
		if strings.Contains(segment, "-") {
			// 按连字符分割为起始和结束部分（最多分割一次，避免IP中包含多个连字符）
			parts := strings.SplitN(segment, "-", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("无效的IP段格式: %s", segment)
			}

			startIPStr := strings.TrimSpace(parts[0]) // 起始IP字符串
			endPart := strings.TrimSpace(parts[1])    // 结束部分（可能是完整IP或最后一段数字）

			// 解析起始IP
			startIP := net.ParseIP(startIPStr)
			if startIP == nil {
				return nil, fmt.Errorf("无效的起始IP: %s", startIPStr)
			}

			// 仅支持IPv4解析
			if ipv4 := startIP.To4(); ipv4 != nil {
				startIP = ipv4 // 确保使用IPv4格式

				// 解析结束部分（可能是完整IP或最后一个八位组的数字）
				var endIP net.IP
				if endIP = net.ParseIP(endPart); endIP != nil {
					// 结束部分是完整IP，转换为IPv4格式
					endIP = endIP.To4()
					if endIP == nil {
						return nil, fmt.Errorf("无效的结束IP: %s", endPart)
					}
				} else {
					// 结束部分不是完整IP，尝试解析为数字（表示IP最后一个八位组）
					endNum, err := strconv.Atoi(endPart)
					if err != nil || endNum < 0 || endNum > 255 {
						return nil, fmt.Errorf("无效的结束部分: %s", endPart)
					}
					// 构造结束IP：复制起始IP的前三个八位组，最后一个八位组替换为解析的数字
					endIP = make(net.IP, len(startIP))
					copy(endIP, startIP)
					endIP[len(endIP)-1] = byte(endNum)
				}

				// 生成从startIP到endIP的所有IP（包含首尾）
				for cmp := bytes.Compare(startIP, endIP); cmp <= 0; cmp = bytes.Compare(startIP, endIP) {
					result = append(result, startIP.String())
					// 递增IP（处理进位，如192.168.1.255 -> 192.168.2.0）
					for i := len(startIP) - 1; i >= 0; i-- {
						startIP[i]++
						if startIP[i] > 0 {
							// 无进位，跳出循环
							break
						}
					}
				}
			} else {
				return nil, fmt.Errorf("IPv6范围解析暂不支持")
			}
		} else {
			// 处理单个IP
			ip := net.ParseIP(segment)
			if ip == nil {
				return nil, fmt.Errorf("无效的IP地址: %s", segment)
			}
			result = append(result, ip.String())
		}
	}

	return result, nil
}
