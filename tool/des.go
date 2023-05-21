package tool

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

const tokenKey = "cc88dd66"

type DesToken struct {
}

func paddingText(str []byte, blockSize int) []byte {
	//需要填充的数据长度
	paddingCount := blockSize - len(str)%blockSize
	//填充数据为：paddingCount ,填充的值为：paddingCount
	paddingStr := bytes.Repeat([]byte{byte(paddingCount)}, paddingCount)
	newPaddingStr := append(str, paddingStr...)
	//fmt.Println(newPaddingStr)
	return newPaddingStr
}

func unPaddingText(str []byte) []byte {
	n := len(str)
	count := int(str[n-1])
	newPaddingText := str[:n-count]
	return newPaddingText
}

func (DesToken) Encrypt(token string) (string, bool) {
	key := []byte(tokenKey)
	plainText := []byte(token)

	//1、创建并返回一个使用DES算法的cipher.Block接口
	block, _ := des.NewCipher(key)
	//2、对数据进行填充
	src1 := paddingText(plainText, block.BlockSize())

	//3.创建一个密码分组为链接模式，底层使用des加密的blockmode接口
	iv := []byte("aaaabbbb")
	blockMode := cipher.NewCBCEncrypter(block, iv)
	//4加密连续的数据块
	desc := make([]byte, len(src1))
	blockMode.CryptBlocks(desc, src1)

	return base64.StdEncoding.EncodeToString(desc), true
}

func (DesToken) Decrypt(token string) (string, bool) {
	// 定义密钥（长度必须是 8 字节）
	key := []byte(tokenKey)

	// 加密后的数据
	cipherText, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", false
	}

	block, _ := des.NewCipher(key)
	iv := []byte("aaaabbbb")
	//链接模式，创建blockMode接口
	blockeMode := cipher.NewCBCDecrypter(block, iv)
	//解密
	blockeMode.CryptBlocks(cipherText, cipherText)
	//去掉填充
	newText := unPaddingText(cipherText)
	return string(newText), true
}
