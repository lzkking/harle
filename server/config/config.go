package config

import (
	"encoding/json"
	"github.com/lzkking/harle/server/assets"
	"os"
	"path/filepath"
)

const (
	ServerConfigPath = "server.json"
	LogFile          = "harle.log"
)

func GetServerConfigPath() string {
	appWorkDir := assets.GetRootAppDir()

	serverConfigPath := filepath.Join(appWorkDir, "configs", ServerConfigPath)

	return serverConfigPath
}

type ServerConfig struct {
	LogFile string `json:"log_file"`
	RunMode string `json:"run_mode"`
}

func (c *ServerConfig) Save() error {
	configPath := GetServerConfigPath()
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0700)
		if err != nil {
			return err
		}
	}
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		panic(err)
	}

	return nil
}

func GetServerConfig() *ServerConfig {
	configPath := GetServerConfigPath()
	config := getDefaultServerConfig()

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return config
		}

		err = json.Unmarshal(data, config)
		if err != nil {
			return config
		}
	}
	err := config.Save()
	if err != nil {
		panic(err)
	}
	return config
}

func getDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		LogFile: filepath.Join(assets.GetRootAppDir(), "log", LogFile),
		RunMode: "DEBUG",
	}
}
