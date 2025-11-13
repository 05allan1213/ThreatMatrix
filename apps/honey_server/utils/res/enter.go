package res

// File: utils/res/response.go
// Description: HTTP 响应统一封装模块

import (
	"honey_server/utils/validate"

	"github.com/gin-gonic/gin"
)

// Response 定义统一的响应结构体
type Response struct {
	Code int    `json:"code"` // 业务状态码（0 表示成功，非 0 表示各种错误）
	Data any    `json:"data"` // 响应数据内容
	Msg  string `json:"msg"`  // 响应说明信息
}

// 底层封装统一输出 JSON 格式响应
func response(code int, data any, msg string, c *gin.Context) {
	c.JSON(200, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// Ok 输出成功响应（可自定义 msg）
func Ok(data any, msg string, c *gin.Context) {
	response(0, data, msg, c)
}

// OkWithData 输出成功响应（默认 msg="成功"）
func OkWithData(data any, c *gin.Context) {
	Ok(data, "成功", c)
}

// OkWithMsg 输出成功响应（仅提示消息，无数据）
func OkWithMsg(msg string, c *gin.Context) {
	Ok(gin.H{}, msg, c)
}

// OkWithList 输出带分页或列表的成功响应
func OkWithList(list any, count int64, c *gin.Context) {
	Ok(gin.H{"list": list, "count": count}, "成功", c)
}

// Fail 输出失败响应（可自定义错误码与消息）
func Fail(code int, msg string, c *gin.Context) {
	response(code, nil, msg, c)
}

// FailWithMsg 输出默认错误码 (1001) 与自定义错误信息
func FailWithMsg(msg string, c *gin.Context) {
	response(1001, nil, msg, c)
}

// FailWithError 接收 error 对象并输出错误信息
func FailWithError(err error, c *gin.Context) {
	msg := validate.ValidateError(err)
	response(1001, nil, msg, c)
}
