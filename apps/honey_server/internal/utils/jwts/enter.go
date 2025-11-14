package jwts

// File: utils/jwts/enter.go
// Description: JWT 生成与解析

import (
	"errors"
	"honey_server/internal/global"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 在 JWT 中携带的用户信息
type ClaimsUserInfo struct {
	UserID uint `json:"userID"` // 用户ID
	Role   int8 `json:"role"`   // 用户角色
}

// Claims 定义自定义的声明体结构，包含：
// - 自定义用户信息 (ClaimsUserInfo)
// - 标准 JWT 声明字段 (jwt.StandardClaims)
type Claims struct {
	ClaimsUserInfo
	jwt.StandardClaims
}

// 根据用户信息生成 JWT
func GetToken(info ClaimsUserInfo) (string, error) {
	j := global.Config.Jwt // 从全局配置中读取 JWT 相关配置

	// 构造 claims
	cla := Claims{
		ClaimsUserInfo: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(j.Expires) * time.Second).Unix(), // 设置过期时间（单位秒）
			Issuer:    j.Issuer,                                                      // 设置签发人
		},
	}

	// 使用 HS256 算法创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cla)

	// 使用配置中的密钥进行签名，生成最终 token 字符串
	return token.SignedString([]byte(j.Secret))
}

// 对传入的 token 字符串进行解析与验证
func ParseToken(tokenString string) (*Claims, error) {
	j := global.Config.Jwt // 获取全局 JWT 配置

	// 使用指定的密钥解析 token，并绑定到自定义的 Claims 结构中
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// 类型断言并校验 token 是否有效
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		// 校验签发人是否匹配，防止伪造
		if claims.Issuer != j.Issuer {
			return nil, errors.New("invalid issuer")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
