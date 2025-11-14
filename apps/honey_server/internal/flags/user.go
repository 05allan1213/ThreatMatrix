package flags

// File: flags/user.go
// Description: 提供通过命令行创建和列出用户的功能

import (
	"encoding/json"
	"fmt"
	"honey_server/internal/global"
	"honey_server/internal/models"
	"honey_server/internal/service/user_service"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

type User struct {
}

// Create 创建用户，可以通过命令行参数或交互式输入方式
func (User) Create(value string) {
	var userInfo user_service.UserCreateRequest

	// 若命令行参数中传入 JSON 字符串，则解析为用户信息
	if value != "" {
		err := json.Unmarshal([]byte(value), &userInfo)
		if err != nil {
			logrus.Errorf("用户信息错误 %s", err)
			return
		}
	} else {
		// 否则，使用命令行交互方式创建用户
		fmt.Println("请选择角色： 1 管理员 2 普通用户")
		_, err := fmt.Scanln(&userInfo.Role)
		if err != nil {
			fmt.Println("输入错误", err)
			return
		}

		// 校验角色输入是否正确
		if !(userInfo.Role == 1 || userInfo.Role == 2) {
			fmt.Println("用户角色输入错误", err)
			return
		}

		// 输入用户名
		fmt.Println("请输入用户名")
		fmt.Scanln(&userInfo.Username)

		// 输入密码（不回显）
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

		// 检查两次密码是否一致
		if string(password) != string(rePassword) {
			fmt.Println("两次密码不一致")
			return
		}

		userInfo.Password = string(password)
	}

	// 使用 user_service 执行创建逻辑
	us := user_service.NewUserService(global.Log)
	_, err := us.Create(userInfo)
	if err != nil {
		logrus.Fatal(err)
	}
}

// List 列出最近创建的10个用户
func (User) List() {
	var userList []models.UserModel

	// 查询最近 10 个用户，按创建时间倒序排列
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
