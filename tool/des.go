package tool

import (
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

const tokenKey = "cc88dd66"

type DesToken struct {
}

func (DesToken) Encrypt(token string) (string, bool) {
	key := []byte(tokenKey)
	plaintext := []byte(token)

	// 创建 DES 加密器
	block, err := des.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return "", false
	}

	// 使用 ECB 模式创建加密器
	mode := cipher.NewCBCEncrypter(block, key)

	// 创建一个切片，用于存储加密后的结果
	ciphertext := make([]byte, len(plaintext))

	// 执行加密操作
	mode.CryptBlocks(ciphertext, plaintext)

	return string(ciphertext), true
}

func (DesToken) Decrypt(token string) (string, bool) {
	// 定义密钥（长度必须是 8 字节）
	key := []byte(tokenKey)

	// 定义待解密的数据
	ciphertext := []byte(token)

	// 创建 DES 解密器
	block, err := des.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating cipher:", err)
		return "", false
	}

	// 使用 ECB 模式创建解密器
	mode := cipher.NewCBCDecrypter(block, key)

	// 创建一个切片，用于存储解密后的结果
	plaintext := make([]byte, len(ciphertext))

	// 执行解密操作
	mode.CryptBlocks(plaintext, ciphertext)

	// 输出解密结果
	return string(plaintext), true
}
