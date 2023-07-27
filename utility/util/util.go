package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

// GenerateSalt 生成salt
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// HashPassword hash密码加密
func HashPassword(password string, salt string) string {
	hashed := sha256.New()
	Bytes, _ := hex.DecodeString(salt)
	_, _ = hashed.Write(Bytes)
	_, _ = hashed.Write([]byte(password))
	return hex.EncodeToString(hashed.Sum(nil))
}

// MergeMap 合并两个map
func MergeMap(m1, m2 map[string]interface{}) map[string]interface{} {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
