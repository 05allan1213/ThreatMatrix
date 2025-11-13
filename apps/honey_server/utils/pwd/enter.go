package pwd

// File: utils/pwd/enter.go
// Description: 密码加密与校验工具

import "golang.org/x/crypto/bcrypt"

// 使用 bcrypt 对明文密码进行加密
func GenerateFromPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// 校验明文密码与哈希是否匹配
func CompareHashAndPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
