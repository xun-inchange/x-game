package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"path/filepath"
	"x-game/server/common/consts"
)

type ClientConfig struct {
	ServerUrl string `toml:"ServerAddr"`
}

var clientConfig *ClientConfig

func init() {
	clientConfig = &ClientConfig{}
	clientConfig.loadData()
}

func (m *ClientConfig) loadData() {
	path := getConfigPath()
	_, err := toml.DecodeFile(path, m)
	if err != nil {
		log.Panicf("load config data err[%v]", err.Error())
	}
}

func getConfigPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, consts.Dir, "client.toml")
}

func GetClientConfig() *ClientConfig {
	return clientConfig
}
