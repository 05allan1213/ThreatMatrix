package captcha

// File: utils/captcha/enter.go
// Description: 图片验证码存储配置。

import "github.com/mojocn/base64Captcha"

var CaptchaStore = base64Captcha.DefaultMemStore // 用来后续校验图片验证码值
