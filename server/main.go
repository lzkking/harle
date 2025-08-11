package main

import (
	"github.com/lzkking/harle/server/internal/transport"
	"github.com/lzkking/harle/server/log"
	"github.com/lzkking/harle/server/pkg/crypto/rsa"
)

func main() {
	log.Init()

	//初始化RSA密钥
	rsa.Init()

	transport.StartHttpListener()
}
