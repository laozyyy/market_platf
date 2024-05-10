package common

import (
	"encoding/base64"
	"math/rand"
)

func GenerateRandomString(length int) string {
	// 计算生成字节数
	byteLength := length * 3 / 4

	// 生成随机字节
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	// 将随机字节编码为字符串
	randomString := base64.URLEncoding.EncodeToString(randomBytes)

	// 截取指定长度的随机字符串
	randomString = randomString[:length]

	return randomString
}
