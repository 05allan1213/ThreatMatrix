package core

// File: core/ip_addr.go
// Description: 使用 ip2region 实现 IP 地址归属地解析的核心组件

import (
	_ "embed"
	"fmt"
	"image_server/internal/utils/ip"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/sirupsen/logrus"
)

var searcher *xdb.Searcher

//go:embed ip2region.xdb
var addrDB []byte

// InitIPDB 初始化 IP 归属地数据库，将嵌入的 xdb 数据加载到内存中
func InitIPDB() {
	// 创建一个默认的 IPv4 版本对象
	version := xdb.IPv4
	_searcher, err := xdb.NewWithBuffer(version, addrDB)
	if err != nil {
		logrus.Fatalf("ip地址数据库加载失败 %s", err)
		return
	}
	searcher = _searcher
}

// GetIpAddr 根据传入的 IP 字符串返回归属地信息
// 优先判断是否为内网 IP，其次使用 ip2region 查询归属地。
// 返回格式根据省市等字段是否为 0 进行拼接组合。
func GetIpAddr(_ip string) (addr string) {
	// 内网地址无需查询数据库
	if ip.HasLocalIPAddr(_ip) {
		return "内网"
	}

	// 从搜索器中查询归属地字符串
	region, err := searcher.SearchByStr(_ip)
	if err != nil {
		logrus.Warnf("错误的ip地址 %s", err)
		return "异常地址"
	}

	// ip2region 默认返回 5 段数据，用"|"分隔
	_addrList := strings.Split(region, "|")
	if len(_addrList) != 5 {
		logrus.Warnf("异常的ip地址 %s", _ip)
		return "未知地址"
	}

	// 五个部分分别为：国家、省份、城市、运营商（以 ip2region 官方格式为准）
	country := _addrList[0]
	province := _addrList[2]
	city := _addrList[3]

	// 优先输出：省份 · 城市
	if province != "0" && city != "0" {
		return fmt.Sprintf("%s·%s", province, city)
	}

	// 其次：国家 · 省份
	if country != "0" && province != "0" {
		return fmt.Sprintf("%s·%s", country, province)
	}

	// 再次：只返回国家字段
	if country != "0" {
		return country
	}

	// 最后，返回 ip2region 原始 Region 字段
	return region
}
