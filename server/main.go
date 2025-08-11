package main

import (
	"github.com/lzkking/harle/server/log"
	"go.uber.org/zap"
)

func main() {
	log.Init()
	zap.S().Infof("server start")
}
