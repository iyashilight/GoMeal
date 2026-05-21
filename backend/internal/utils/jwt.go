package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 自定义声明，嵌入标准注册声明
type Claims struct {
	UserID   uint   `json:"user_id"`   // 用户 ID
	Phone    string `json:"phone"`     // 手机号
	UserType int    `json:"user_type"` // 用户类型
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token，使用 HS256 签名
// 参数：用户信息 + 密钥 + 过期时间
func GenerateToken(userID uint, phone string, userType int, secret string,
	expireHours int) (string, error) {
	claims := Claims{
		UserID:   userID,
		Phone:    phone,
		UserType: userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour *
				time.Duration(expireHours))),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并验证 JWT Token，返回声明信息
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
