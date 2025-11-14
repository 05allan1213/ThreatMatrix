package user_api

// File:api/user_api/login.go
// Description: 用户登录接口

import (
	"honey_server/internal/global"
	"honey_server/internal/middleware"
	"honey_server/internal/models"
	"honey_server/internal/service/log_service"
	"honey_server/internal/utils/captcha"
	"honey_server/internal/utils/jwts"
	"honey_server/internal/utils/pwd"
	"honey_server/internal/utils/res"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username    string `json:"username" binding:"required" label:"用户名"`
	Password    string `json:"password" binding:"required" label:"密码"`
	CaptchaID   string `json:"captchaID" binding:"required" label:"验证码ID"`
	CaptchaCode string `json:"captchaCode" binding:"required" label:"验证码"`
}

// 用户登录接口
func (UserApi) LoginView(c *gin.Context) {
	// 从请求中绑定并获取登录参数（如用户名、密码、验证码）
	cr := middleware.GetBind[LoginRequest](c)
	// log := middleware.GetLog(c) // 若需要日志可启用
	loginLog := log_service.NewLoginLog(c) // 创建登录日志服务实例

	// 校验验证码是否为空
	if cr.CaptchaID == "" || cr.CaptchaCode == "" {
		loginLog.FailLog(cr.Username, "", "未输入图片验证码")
		res.FailWithMsg("请输入图片验证码", c)
		return
	}

	// 验证图片验证码是否正确（第三个参数true表示验证后删除验证码）
	if !captcha.CaptchaStore.Verify(cr.CaptchaID, cr.CaptchaCode, true) {
		loginLog.FailLog(cr.Username, "", "图片验证码验证失败")
		res.FailWithMsg("图片验证码验证失败", c)
		return
	}

	// 从数据库中查询用户信息
	var user models.UserModel
	err := global.DB.Take(&user, "username = ?", cr.Username).Error
	if err != nil {
		// 用户不存在或查询出错，返回错误提示
		loginLog.FailLog(cr.Username, cr.Password, "用户名不存在")
		res.FailWithMsg("用户名或密码错误", c)
		return
	}

	// 校验用户密码是否匹配
	if !pwd.CompareHashAndPassword(user.Password, cr.Password) {
		loginLog.FailLog(cr.Username, cr.Password, "密码错误")
		res.FailWithMsg("用户名或密码错误", c)
		return
	}

	// 生成 JWT Token，包含用户ID和角色信息
	token, err := jwts.GetToken(jwts.ClaimsUserInfo{
		UserID: user.ID,
		Role:   user.Role,
	})
	if err != nil {
		// Token 生成失败记录日志并返回错误
		logrus.Errorf("生成token失败 %s", err)
		res.FailWithMsg("登录失败", c)
		return
	}

	now := time.Now().Format(time.DateTime)
	global.DB.Model(&user).Update("last_login_date", now) // 更新最后登录时间

	// 登录成功，返回生成的 Token，并记录登录成功日志
	loginLog.SuccessLog(user.ID, cr.Username)
	res.OkWithData(token, c)
}
