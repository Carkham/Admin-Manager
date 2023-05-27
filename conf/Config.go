package conf

import (
	"gopkg.in/ini.v1"
	"path/filepath"
	"runtime"
)

const iniPath = "conf/conf.ini"

var Config GlobalConfigDefine

func init() {
	// This is not work well for 'go run' but for 'go build'
	_, path, _, ok := runtime.Caller(0)
	if ok != true {
		panic("Can't get current file path")
	}
	iniPath := filepath.Join(filepath.Dir(path), filepath.Base(iniPath))
	cfg, err := ini.Load(iniPath)
	if err != nil {
		panic(err)
	}
	err = cfg.MapTo(&Config)
	if err != nil {
		panic(err)
	}
}
