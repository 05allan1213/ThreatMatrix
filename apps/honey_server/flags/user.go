package flags

import (
	"encoding/json"
	"fmt"
	"honey_server/global"
	"honey_server/models"
	"honey_server/utils/pwd"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type User struct {
}

// 用户创建请求结构体，用于接收JSON或交互式输入
type UserInfoRequest struct {
	Role     int8   `json:"role"`     // 用户角色：1管理员 2普通用户
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// Create 创建用户，可以通过命令行参数或交互式输入方式
func (User) Create(value string) {
	var userInfo UserInfoRequest

	// 如果传入了 JSON 参数，则从 JSON 解析用户信息
	if value != "" {
		err := json.Unmarshal([]byte(value), &userInfo)
		if err != nil {
			logrus.Errorf("用户信息错误 %s", err)
			return
		}
	} else {
		// 未传入 JSON，则使用交互式输入
		fmt.Println("请选择角色： 1 管理员 2 普通用户")
		_, err := fmt.Scanln(&userInfo.Role)
		if err != nil {
			fmt.Println("输入错误", err)
			return
		}
		// 校验角色输入是否合法
		if !(userInfo.Role == 1 || userInfo.Role == 2) {
			fmt.Println("用户角色输入错误", err)
			return
		}

		// 输入用户名
		fmt.Println("请输入用户名")
		fmt.Scanln(&userInfo.Username)

		// 输入密码（隐藏输入）
		fmt.Println("请输入密码")
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("读取密码时出错:", err)
			return
		}

		// 确认密码
		fmt.Println("请再次输入密码")
		rePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println("读取密码时出错:", err)
			return
		}

		// 比对两次输入的密码是否一致
		if string(password) != string(rePassword) {
			fmt.Println("两次密码不一致")
			return
		}

		userInfo.Password = string(password)
	}

	// 检查用户名是否已存在
	var u models.UserModel
	err := global.DB.Take(&u, "username = ?", userInfo.Username).Error
	if err == nil {
		fmt.Println("用户名已存在")
		return
	}

	// 对密码进行加密
	hashPwd, _ := pwd.GenerateFromPassword(userInfo.Password)

	// 创建用户记录
	err = global.DB.Create(&models.UserModel{
		Username: userInfo.Username,
		Password: hashPwd,
		Role:     userInfo.Role,
	}).Error
	if err != nil {
		logrus.Errorf("用户创建失败 %s", err)
		return
	}

	// 输出成功信息
	logrus.Infof("用户创建成功")
}

// List 列出最近创建的10个用户
func (User) List() {
	var userList []models.UserModel

	// 查询最近10个用户（按创建时间倒序）
	global.DB.Order("created_at desc").Limit(10).Find(&userList)

	// 输出用户信息
	for _, model := range userList {
		fmt.Printf("用户id：%d  用户名：%s 用户角色：%d 创建时间：%s\n",
			model.ID,
			model.Username,
			model.Role,
			model.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}
}
