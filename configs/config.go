package configs

import (
	"path/filepath"
	"twimg/utils"
)

var mode int = 0
var configPath = "configs"

// Deployment 部署模式 0: debug 1: release
func Deployment() bool {
	if mode == 0 {
		return true
	}
	return false
}

// readFile 读取 JSON 配置文件
func readFile(name string) (d map[string]interface{}) {
	cfgPath := filepath.Join(utils.FileSuite.LocalPath(Deployment()), configPath, name)
	if utils.FileSuite.CheckExist(cfgPath) {
		d = utils.DataSuite.RawMap2Map(utils.FileSuite.Read(cfgPath))
		return d
	}
	return
}

// APIKeys Twitter API Keys
func APIKeys() (d map[string]interface{}) {
	d = readFile("apikeys.json")
	return
}
