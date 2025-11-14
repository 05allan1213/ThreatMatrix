package docker

// File: utils/docker/parse_manifest.go
// Description: 提供对本地 Docker 镜像文件（tar/gz）的解析能力，从 manifest.json 中提取镜像名称、Tag、ImageID。

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

const manifestFile = "manifest.json"

// ParseImageMetadata 解析本地 Docker 镜像导出文件的基础元信息（ImageID、镜像名、Tag）。
func ParseImageMetadata(filePath string) (string, string, string, error) {
	file, err := os.Open(filePath) // 打开镜像文件
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	var reader io.Reader = file

	// 若文件后缀为 .gz，则创建 gzip.Reader 进行解压，否则保持 tarReader 使用原始 file。
	if strings.HasSuffix(filePath, ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return "", "", "", err
		}
		defer gzReader.Close()
		reader = gzReader
	}

	// tar.Reader 用于遍历镜像文件内部目录结构
	tarReader := tar.NewReader(reader)

	var imageID, imageName, imageTag string

	for {
		header, err := tarReader.Next() // 读取下一个 tar header
		if err == io.EOF {
			break // 文件读完
		}
		if err != nil {
			return "", "", "", err
		}

		switch header.Name {
		case manifestFile: // 匹配 manifest.json（Docker 镜像导出文件的核心描述文件）
			manifestData, err := io.ReadAll(tarReader) // 读取 manifest.json 内容
			if err != nil {
				return "", "", "", err
			}

			// 解析 manifest.json
			data, err := extractImage(string(manifestData))
			if err != nil {
				return "", "", "", err
			}
			return data.ImageID, data.ImageName, data.ImageTag, nil
		}
	}

	// 若未正确提取到信息，返回错误
	if imageID == "" || imageName == "" || imageTag == "" {
		return "", "", "", fmt.Errorf("无法从镜像文件中提取完整的元数据")
	}

	return imageID, imageName, imageTag, nil
}

// manifestType 对应 manifest.json 的结构
type manifestType struct {
	Config   string   `json:"Config"`   // 指向镜像 config 文件，路径类似 blobs/sha256/xxxx.json
	RepoTags []string `json:"RepoTags"` // 镜像名与 Tag 列表，如 ["nginx:1.21"]
}

// manifestData 封装提取后的镜像元数据（名称、Tag、ImageID）
type manifestData struct {
	ImageID   string
	ImageName string
	ImageTag  string
}

// extractImage 解析 manifest.json 内容获取镜像名称、Tag、ImageID
func extractImage(manifest string) (data manifestData, err error) {
	var t []manifestType

	// 使用 JSON 反序列化解析 manifest 内容
	err = json.Unmarshal([]byte(manifest), &t)
	if err != nil {
		err = fmt.Errorf("解析manifest文件失败 %s", err)
		return
	}

	if len(t) == 0 {
		err = fmt.Errorf("解析manifest文件内容失败 %s", manifest)
		return
	}

	// RepoTags 示例：["nginx:1.21"]
	if len(t[0].RepoTags) == 0 {
		err = fmt.Errorf("manifest文件中没有找到RepoTags信息")
		return
	}
	
	repoTags := t[0].RepoTags[0]
	_list := strings.Split(repoTags, ":")

	// 镜像名称与 Tag
	if len(_list) < 2 {
		err = fmt.Errorf("无效的RepoTags格式: %s", repoTags)
		return
	}
	data.ImageName = _list[0]
	data.ImageTag = _list[1]

	// Config 示例 "blobs/sha256/xxxxxxxxxxxxxx.json"
	// 通常镜像 ID 使用其中 sha256 的前 12 位
	configParts := strings.Split(t[0].Config, "/")
	if len(configParts) < 3 {
		err = fmt.Errorf("无效的Config格式: %s", t[0].Config)
		return
	}
	
	configFileName := configParts[2]
	if len(configFileName) < 12 {
		err = fmt.Errorf("Config文件名长度不足12位: %s", configFileName)
		return
	}
	
	data.ImageID = configFileName[:12]

	return
}