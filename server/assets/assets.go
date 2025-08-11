package assets

import (
	"os"
	"path/filepath"
)

const (
	ServerWorkDirEnv = "ServerWorkDirEnv"
)

// GetRootAppDir - 获取工作路径
func GetRootAppDir() string {
	value := os.Getenv(ServerWorkDirEnv)
	var dir string

	if len(value) == 0 {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dir = filepath.Join(wd, ".harle")
	} else {
		dir = value
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			panic(err)
		}
	}
	return dir
}
