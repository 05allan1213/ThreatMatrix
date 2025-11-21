package core

// File: core/manf.go
// Description: 通过MAC地址查询对应的厂商信息

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// OUIDatabase 用于存储OUI（组织唯一标识符）到厂商名称的映射关系
type OUIDatabase struct {
	vendors map[string]string // 键：标准化的OUI（6位十六进制字符串），值：厂商名称
}

// NewOUIDatabase 创建一个新的OUI数据库实例
func NewOUIDatabase() *OUIDatabase {
	return &OUIDatabase{
		vendors: make(map[string]string), // 初始化映射表
	}
}

// LoadFromIEEE 从IEEE标准格式的oui.txt文件读取器加载OUI数据
func (db *OUIDatabase) LoadFromIEEE(reader *bufio.Reader) error {
	// 正则表达式：匹配oui.txt中的OUI记录行（格式示例：00-1B-44   (hex)  Intel Corporation）
	// 分组1：OUI部分（如00-1B-44），分组2：厂商名称（如Intel Corporation）
	re := regexp.MustCompile(`^([0-9A-Fa-f]{2}[-:][0-9A-Fa-f]{2}[-:][0-9A-Fa-f]{2})\s+\(hex\)\s+(.*)$`)

	scanner := bufio.NewScanner(reader)
	lineNum := 0 // 记录行号，用于错误定位（未直接使用，保留扩展空间）

	// 逐行扫描文件内容
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // 跳过空行
		}

		// 尝试匹配OUI记录行
		matches := re.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue // 忽略非记录行（如注释、标题、空行等）
		}

		// 提取OUI并标准化：去除分隔符（-或:），转为大写
		ouiRaw := matches[1]
		oui := strings.ToUpper(strings.ReplaceAll(ouiRaw, "-", ""))
		oui = strings.ReplaceAll(oui, ":", "") // 处理冒号分隔的格式（如00:1B:44）

		// 提取厂商名称并去除首尾多余空格
		vendor := strings.TrimSpace(matches[2])

		// 将标准化的OUI和厂商名称存入映射表
		db.vendors[oui] = vendor
	}

	// 检查扫描过程中是否出现错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("扫描文件错误: %v", err)
	}

	fmt.Printf("成功加载OUI数据库，共加载 %d 条记录\n", len(db.vendors))
	return nil
}

// LookupVendor 通过MAC地址查询对应的厂商名称
func (db *OUIDatabase) LookupVendor(mac string) (string, bool) {
	// 标准化MAC地址：去除所有分隔符（:/-/.），转为大写
	mac = strings.ToUpper(
		strings.ReplaceAll(
			strings.ReplaceAll(
				strings.ReplaceAll(mac, ":", ""),
				"-", ""),
			".", ""),
	)

	// 提取MAC地址的前6位（OUI部分），若长度不足则返回未找到
	if len(mac) < 6 {
		return "", false
	}
	oui := mac[:6]

	// 从映射表中查询厂商名称
	vendor, exists := db.vendors[oui]
	return vendor, exists
}

// oui 嵌入式的OUI数据文件（oui.txt）内容，通过go:embed指令嵌入到二进制文件中
//
//go:embed oui.txt
var oui []byte

var manufDB *OUIDatabase // 全局的OUI数据库实例

// init 初始化函数：创建并加载全局OUI数据库
func init() {
	manufDB = NewOUIDatabase()
	// 从嵌入式的oui数据创建读取器，并加载到数据库
	err := manufDB.LoadFromIEEE(bufio.NewReader(bytes.NewReader(oui)))
	if err != nil {
		logrus.Fatalf("加载OUI数据库失败: %v", err)
		return
	}
}

// ManufQuery 对外提供的MAC地址厂商查询接口
func ManufQuery(mac string) (string, bool) {
	return manufDB.LookupVendor(mac)
}
