package ecdh

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/hkdf"
	"io"
)

// GetTmpKey - 获取临时的加密的对称密钥
func GetTmpKey() (publicKey []byte, privateKey []byte, err error) {
	curve := ecdh.X25519()
	private, err := curve.GenerateKey(rand.Reader)
	if err != nil {
		return
	}
	privateKey = private.Bytes()
	publicKey = private.PublicKey().Bytes()
	return
}

// CalcSessionKey - 计算会话密钥
func CalcSessionKey(serverPublicKey, agentPrivateKey, agentPublicKey []byte) (keySend []byte, keyRecv []byte, err error) {
	curve := ecdh.X25519()
	private, err := curve.NewPrivateKey(agentPrivateKey)
	if err != nil {
		return
	}

	public, err := curve.NewPublicKey(serverPublicKey)
	if err != nil {
		return
	}

	shared, err := private.ECDH(public)
	if err != nil {
		return
	}

	//	生成会话密钥
	salt := sha256.Sum256(append(serverPublicKey, agentPublicKey...))
	keySend = hkdfBytes(shared, salt[:], []byte("B->A key v1"), 32)
	keyRecv = hkdfBytes(shared, salt[:], []byte("A->B key v1"), 32)
	return
}

func hkdfBytes(secret, salt, info []byte, n int) []byte {
	h := hkdf.New(sha256.New, secret, salt, info)
	out := make([]byte, n)
	if _, err := io.ReadFull(h, out); err != nil {
		panic(err)
	}
	return out
}
