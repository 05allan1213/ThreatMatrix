package captcha_api

// File: api/captcha_api/captcha.go
// Description: 图片验证码接口，用于生成前端展示的验证码图片。

import (
	"honey_server/utils/captcha"
	"honey_server/utils/res"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/sirupsen/logrus"
)

// CaptchaApi 定义验证码接口结构体（空结构，仅用于方法接收者）
type CaptchaApi struct{}

// GenerateResponse 定义生成验证码接口的响应结构
type GenerateResponse struct {
	CaptchaID string `json:"captchaID"` // 验证码唯一标识符，用于验证时关联验证码内容
	Captcha   string `json:"captcha"`   // 验证码图片的 base64 编码字符串，前端可直接渲染显示
}

// GenerateView 生成验证码图片
// 1. 初始化验证码参数（宽高、噪点、长度、字符源等）
// 2. 使用 base64Captcha 生成图片和 ID
// 3. 若生成失败则返回错误响应，否则返回 base64 图片与对应 ID
func (CaptchaApi) GenerateView(c *gin.Context) {
	// 定义验证码驱动配置
	var driver = base64Captcha.DriverString{
		Width:           200,          // 图片宽度
		Height:          60,           // 图片高度
		NoiseCount:      2,            // 干扰噪点数量
		ShowLineOptions: 4,            // 干扰线样式选项
		Length:          4,            // 验证码长度
		Source:          "0123456789", // 验证码字符来源，仅数字
	}

	// 创建验证码对象并指定全局存储
	cp := base64Captcha.NewCaptcha(&driver, captcha.CaptchaStore)

	// 生成验证码：返回 id（验证码标识）、b64s（图片内容）
	id, b64s, _, err := cp.Generate()
	if err != nil {
		// 生成失败时记录日志并返回错误响应
		logrus.Errorf("图片验证码生成失败 %s", err)
		res.FailWithMsg("图片验证码生成失败", c)
		return
	}

	// 成功返回：包含验证码ID与base64编码图片
	res.OkWithData(GenerateResponse{
		CaptchaID: id,
		Captcha:   b64s,
	}, c)
}
