package common

type Config struct {
	Domain     string `required:"true"`
	AppName    string `required:"true"`
	BaseUrl    string `required:"true"`
	Port       int    `required:"true"`
	Template   string `required:"true"`
	StaticPath string `required:"true"`
	DBPath     string `required:"true"`
	LogPath    string `required:"true"`
	Debug      bool
	EnableWeb  bool
	MySQL      MySQLConfig
	Redis      RedisConfig
}

type MySQLConfig struct {
	Host   string `required:"true"`
	User   string `required:"true"`
	Passwd string `required:"true"`
	DB     string `required:"true"`
}

type RedisConfig struct {
	Master string `required:"true"`
	Slave  string
}
