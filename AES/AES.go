package AES

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
)

// 原文,
// AES 加密密钥（32 字节，256 位）
func Encrypt(plaintext []byte, key []byte) []byte {

	// 创建 AES 分组，使用密钥
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("Error creating AES cipher:", err)
		return nil
	}

	// 初始化向量 IV（必须是 16 字节）
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("Error generating IV:", err)
		return nil
	}

	// 创建加密器
	encrypter := cipher.NewCFBEncrypter(block, iv)

	// 创建加密后的数据缓冲区
	ciphertext := make([]byte, len(plaintext))

	// 加密数据
	encrypter.XORKeyStream(ciphertext, plaintext)

	// 将 IV 附加到加密后的数据中
	ciphertext = append(iv, ciphertext...)

	return ciphertext
}

// Decrypt 使用AES解密给定的加密数据
func Decrypt(ciphertext []byte, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 创建解密器
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	decrypter.XORKeyStream(decrypted, ciphertext)

	return decrypted
}
