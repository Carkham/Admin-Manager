package conf

import "gopkg.in/ini.v1"

const iniPath = "conf/conf.ini"

var Config GlobalConfigDefine

func init() {
	cfg, err := ini.Load(iniPath)
	if err != nil {
		panic(err)
	}
	err = cfg.MapTo(&Config)
	if err != nil {
		panic(err)
	}
}
