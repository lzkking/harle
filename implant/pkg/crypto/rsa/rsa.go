package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lzkking/harle/implant/assets"
)

func RsaEncryptData(plainData []byte, publicKey *rsa.PublicKey) (encryptedData []byte, err error) {
	encryptedData, err = rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		plainData,
		nil,
	)
	return
}

func LoadRsaPublicKey() (*rsa.PublicKey, error) {
	pubKeyData := assets.RsaPublicPem
	block, _ := pem.Decode(pubKeyData)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("公钥文件格式错误")
	}

	// 解析为接口
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %v", err)
	}

	// 断言为 *rsa.PublicKey
	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥类型不是 RSA")
	}

	return pubKey, nil
}
