package agent

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/hkdf"
	"io"
	"time"
)

type Agent struct {
	AgentId          string    // 被控端的唯一标识符号
	sessionKey       []byte    // 会话密钥 AES
	sessionKeyRecv   []byte    // 接收数据后的解密密钥
	sessionKeySend   []byte    // 发送数据的加密密钥
	sessionKeyExpire time.Time // 会话密钥过期时间
	serverPublicKey  []byte    // 服务端的临时公钥,针对单个agent的,不同agent不同
	serverPrivateKey []byte    // 服务端的临时私钥,针对单个agent的,是agent临时公钥的密钥对
	AgentPublicKey   []byte    //agent传递来的临时公钥
}

// GetSessionKey - 获取会话密钥
func (a *Agent) GetSessionKey() []byte {
	return a.sessionKey
}

// CalcSessionKey - 计算会话密钥
func (a *Agent) CalcSessionKey() (err error) {
	if len(a.serverPrivateKey) == 0 ||
		len(a.serverPublicKey) == 0 ||
		len(a.AgentPublicKey) == 0 {
		return fmt.Errorf("公私钥目前并未正确传递")
	}

	curve := ecdh.X25519()
	private, err := curve.NewPrivateKey(a.serverPrivateKey)
	if err != nil {
		return
	}

	public, err := curve.NewPublicKey(a.AgentPublicKey)
	if err != nil {
		return
	}

	shared, err := private.ECDH(public) //原始共享密钥
	if err != nil {
		return
	}

	a.sessionKey = shared

	//	生成会话密钥
	salt := sha256.Sum256(append(a.serverPublicKey, a.AgentPublicKey...))
	keySend := hkdfBytes(shared, salt[:], []byte("A->B key v1"), 32)
	KeyRecv := hkdfBytes(shared, salt[:], []byte("B->A key v1"), 32)

	a.sessionKeySend = keySend
	a.sessionKeyRecv = KeyRecv
	return
}

// UpdateSessionKey - 更新会话密钥
func (a *Agent) UpdateSessionKey(sessionKey []byte) {
	a.sessionKey = sessionKey
}

// GetServerPublicKey - 获取服务端的临时公钥
func (a *Agent) GetServerPublicKey() []byte {
	if len(a.serverPublicKey) == 0 {
		curve := ecdh.X25519()
		private, err := curve.GenerateKey(rand.Reader)
		if err != nil {
			return nil
		}
		a.serverPublicKey = private.PublicKey().Bytes()
		a.serverPrivateKey = private.Bytes()
	}
	return a.serverPublicKey
}

// GetServerPrivateKey - 获取服务端临时私钥
func (a *Agent) GetServerPrivateKey() []byte {
	if len(a.serverPrivateKey) == 0 {
		curve := ecdh.X25519()
		private, err := curve.GenerateKey(rand.Reader)
		if err != nil {
			return nil
		}
		a.serverPublicKey = private.PublicKey().Bytes()
		a.serverPrivateKey = private.Bytes()
	}
	return a.serverPrivateKey
}

// GetServerTempKey - 获取临时密钥对:公钥、私钥
func (a *Agent) GetServerTempKey() (serverPublicKey []byte, serverPrivateKey []byte) {
	if len(a.serverPrivateKey) == 0 || len(a.serverPublicKey) == 0 {
		curve := ecdh.X25519()
		private, err := curve.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil
		}
		a.serverPublicKey = private.PublicKey().Bytes()
		a.serverPrivateKey = private.Bytes()
	}
	serverPublicKey = a.serverPublicKey
	serverPrivateKey = a.serverPrivateKey
	return
}

func (a *Agent) EncryptData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("data长度为0")
	}

	if len(a.sessionKeySend) == 0 {
		err := a.CalcSessionKey()
		if err != nil {
			return nil, fmt.Errorf("获取加密密钥失败")
		}
	}

	sessionKey := a.sessionKeySend
	block, err := aes.NewCipher(sessionKey)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aead.Seal(nil, nonce, data, nil)

	out := append(nonce, ciphertext...)

	return out, nil
}

func hkdfBytes(secret, salt, info []byte, n int) []byte {
	h := hkdf.New(sha256.New, secret, salt, info)
	out := make([]byte, n)
	if _, err := io.ReadFull(h, out); err != nil {
		panic(err)
	}
	return out
}
