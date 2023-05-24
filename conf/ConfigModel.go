package conf

type GlobalConfigDefine struct {
	Service ServiceConfigDefine `ini:"Service"`
	K8S     K8SConfigDefine     `ini:"K8S"`
	MySQL   MySQLConfigDefine   `ini:"MySQL"`
	Redis   RedisConfigDefine   `ini:"Redis"`
}

type ServiceConfigDefine struct {
	HttpPort       int    `ini:"HttpPort"`
	DeploymentAddr string `ini:"DeploymentAddr"`
}

type K8SConfigDefine struct {
	Address    string `ini:"Address"`
	ConfigPath string `ini:"ConfigPath"`
}

type MySQLConfigDefine struct {
	Address  string `ini:"Address"`
	Port     int    `ini:"Port"`
	Username string `ini:"Username"`
	Password string `ini:"Password"`
	DBName   string `ini:"DBName"`
}

type RedisConfigDefine struct {
	Address  string `ini:"Address"`
	Port     int    `ini:"Port"`
	Password string `ini:"Password"`
}
