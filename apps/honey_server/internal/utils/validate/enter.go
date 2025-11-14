package validate

// File: utils/validate/enter.go
// Description: 参数验证器初始化与错误翻译工具

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// trans 是全局翻译器，用于将校验错误翻译成中文
var trans ut.Translator

// init 函数在包初始化时运行，设置验证器与翻译器
func init() {
	// 创建中文翻译器实例
	uni := ut.New(zh.New())
	trans, _ = uni.GetTranslator("zh")

	// 获取 Gin 默认的验证引擎（validator）
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		// 注册默认的中文翻译
		_ = zh_translations.RegisterDefaultTranslations(v, trans)
	}

	// 自定义字段名映射函数
	// 若结构体字段包含标签 `label:"用户名"`，则错误提示中显示该中文标签
	// 否则默认使用字段名（如 Username）
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
}

// 将验证错误转换为可读的中文字符串
func ValidateError(err error) string {
	// 类型断言，判断错误是否为验证错误集合
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 若不是验证错误，直接返回原始错误信息
		return err.Error()
	}

	var list []string
	// 遍历每个字段错误并翻译成中文
	for _, e := range errs {
		list = append(list, e.Translate(trans))
	}

	// 将所有错误用分号拼接成一个字符串
	return strings.Join(list, ";")
}
