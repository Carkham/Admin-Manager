package conf

type GlobalConfigDefine struct {
	Service    ServiceConfigDefine    `ini:"Service"`
	K8S        K8SConfigDefine        `ini:"K8S"`
	MySQL      MySQLConfigDefine      `ini:"MySQL"`
	Redis      RedisConfigDefine      `ini:"Redis"`
	Webshell   WebShellConfigDefine   `ini:"WebShell"`
	FileSystem FileSystemConfigDefine `ini:"FileSystem"`
	Admin      AdminConfigDefine      `ini:"Admin"`
}

type ServiceConfigDefine struct {
	HttpPort       int    `ini:"HttpPort"`
	RpcPort        string `ini:"RpcPort"`
	DeploymentAddr string `ini:"DeploymentAddr"`
	ExposeIP       string `ini:"ExposeIP"`
	TimeZone       string `ini:"TimeZone"`
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

type WebShellConfigDefine struct {
	TimeoutSec int    `ini:"TimeoutSecond"`
	MinPort    int    `ini:"MinPort"`
	MaxPort    int    `ini:"MaxPort"`
	GoTTYExec  string `ini:"GottyExec"`
	PasswdLen  int    `ini:"PasswordLen"`
	BindIP     string `ini:"BindIP"`
}

type FileSystemConfigDefine struct {
	RootPath string `ini:"RootPath"`
	NFSAddr  string `ini:"NFSAddr"`
}

type AdminConfigDefine struct {
	Username string `ini:"Username"`
	Password string `ini:"Password"`
	Email    string `ini:"Email"`
}
