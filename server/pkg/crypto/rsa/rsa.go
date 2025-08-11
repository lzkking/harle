package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/lzkking/harle/server/config"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

var (
	RsaPrivateKey *rsa.PrivateKey
	RsaPublicKey  *rsa.PublicKey
)

func GenerateRsaKey() error {
	globalConfig := config.GetServerConfig()
	rsaPrivateFile := globalConfig.RsaPrivateKeyFile
	rsaPublicFile := globalConfig.RsaPublicKeyFile

	dir := filepath.Dir(rsaPrivateFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(rsaPrivateFile); err == nil {
		_ = os.Remove(rsaPrivateFile)
	}

	if _, err := os.Stat(rsaPublicFile); err == nil {
		_ = os.Remove(rsaPublicFile)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateFile, err := os.Create(rsaPrivateFile)
	if err != nil {
		return err
	}
	defer privateFile.Close()

	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	err = pem.Encode(privateFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	})
	if err != nil {
		return err
	}

	publicKey := &privateKey.PublicKey
	publicFile, err := os.Create(rsaPublicFile)
	if err != nil {
		return err
	}
	defer publicFile.Close()

	publicBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	err = pem.Encode(publicFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})

	if err != nil {
		return err
	}

	return nil
}

func LoadRSAPrivateKeyFromPEM(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY": // PKCS#1 格式
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY": // PKCS#8 格式
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if k, ok := key.(*rsa.PrivateKey); ok {
			return k, nil
		}
		return nil, errors.New("not RSA private key")
	default:
		return nil, fmt.Errorf("unsupported private key type: %s", block.Type)
	}
}

func LoadRSAPublicKeyFromPEM(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}

	switch block.Type {
	case "PUBLIC KEY": // PKIX 格式
		ifc, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if k, ok := ifc.(*rsa.PublicKey); ok {
			return k, nil
		}
		return nil, errors.New("not RSA public key")
	case "RSA PUBLIC KEY": // PKCS#1 格式
		return x509.ParsePKCS1PublicKey(block.Bytes)
	default:
		return nil, fmt.Errorf("unsupported public key type: %s", block.Type)
	}
}

func Init() {
	globalConfig := config.GetServerConfig()
	rsaPrivateFile := globalConfig.RsaPrivateKeyFile
	rsaPublicFile := globalConfig.RsaPublicKeyFile

	if _, err := os.Stat(rsaPrivateFile); err != nil {
		err = GenerateRsaKey()
		if err != nil {
			zap.S().Errorf("生成RSA密钥失败")
			panic(err)
		}
	}

	var err error
	RsaPrivateKey, err = LoadRSAPrivateKeyFromPEM(rsaPrivateFile)
	if err != nil {
		zap.S().Errorf("加载RSA私钥失败\n")
		panic(err)
	}

	RsaPublicKey, err = LoadRSAPublicKeyFromPEM(rsaPublicFile)
	if err != nil {
		zap.S().Errorf("加载RSA公钥失败\n")
		panic(err)
	}

}

func EncryptWithPublicKey(plaintext []byte) ([]byte, error) {
	h := sha256.New()
	return rsa.EncryptOAEP(h, rand.Reader, RsaPublicKey, plaintext, nil)
}

func DecryptWithPrivateKey(ciphertext []byte) ([]byte, error) {
	h := sha256.New()
	return rsa.DecryptOAEP(h, rand.Reader, RsaPrivateKey, ciphertext, nil)
}

// -------------------- 私钥签名 / 公钥验签（PSS） --------------------

func SignWithPrivateKey(message []byte) ([]byte, error) {
	sum := sha256.Sum256(message)
	// 可根据需要自定义 PSS 选项；nil 等价于默认（SaltLengthAuto）
	return rsa.SignPSS(rand.Reader, RsaPrivateKey, crypto.SHA256, sum[:], nil)
}

func VerifyWithPublicKey(message, signature []byte) error {
	sum := sha256.Sum256(message)
	return rsa.VerifyPSS(RsaPublicKey, crypto.SHA256, sum[:], signature, nil)
}
