package flags

// File: flags/user.go
// Description: 实现用户相关的命令行参数处理逻辑，包括创建用户和列出用户。

import (
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

// Create 创建用户
func (User) Create() {
	var user models.UserModel

	// 选择角色
	fmt.Println("请选择角色： 1 管理员 2 普通用户")
	_, err := fmt.Scanln(&user.Role)
	if err != nil {
		fmt.Println("输入错误", err)
		return
	}
	// 检查角色输入是否合法
	if !(user.Role == 1 || user.Role == 2) {
		fmt.Println("用户角色输入错误", err)
		return
	}

	// 输入用户名
	fmt.Println("请输入用户名")
	fmt.Scanln(&user.Username)

	// 检查用户名是否已存在
	var u models.UserModel
	err = global.DB.Take(&u, "username = ?", user.Username).Error
	if err == nil {
		fmt.Println("用户名已存在")
		return
	}

	// 输入密码（隐藏输入）
	fmt.Println("请输入密码")
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("读取密码时出错:", err)
		return
	}

	// 再次确认密码
	fmt.Println("请再次输入密码")
	rePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("读取密码时出错:", err)
		return
	}

	// 比较两次输入的密码是否一致
	if string(password) != string(rePassword) {
		fmt.Println("两次密码不一致")
		return
	}

	// 生成加密后的密码
	hashPwd, _ := pwd.GenerateFromPassword(string(password))

	// 将用户信息写入数据库
	err = global.DB.Create(&models.UserModel{
		Username: user.Username,
		Password: hashPwd,
		Role:     user.Role,
	}).Error
	if err != nil {
		logrus.Errorf("用户创建失败 %s", err)
		return
	}

	// 创建成功日志输出
	logrus.Infof("用户创建成功")
}

// List 查看最近创建的用户列表（最多显示10条）
func (User) List() {
	var userList []models.UserModel

	// 按创建时间倒序查询最近10个用户
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
