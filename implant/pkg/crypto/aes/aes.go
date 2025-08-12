package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func DecryptData(sessionKey, in []byte) ([]byte, error) {
	if len(in) == 0 {
		return nil, fmt.Errorf("ciphertext 为空")
	}

	// 创建 AES 分组器
	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		return nil, err
	}

	// GCM 封装
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 拆出 nonce 和密文
	nonceSize := aead.NonceSize()
	if len(in) < nonceSize {
		return nil, fmt.Errorf("密文长度不足")
	}

	nonce := in[:nonceSize]      // 前 nonceSize 字节是随机数
	ciphertext := in[nonceSize:] // 后面是密文 + AuthTag

	// 解密并验证
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("解密失败或数据被篡改: %w", err)
	}

	return plaintext, nil
}
