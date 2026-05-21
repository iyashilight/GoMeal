package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword 对明文密码进行 bcrypt 哈希，返回哈希后的密码字符串
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证明文密码是否与 bcrypt 哈希匹配
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
