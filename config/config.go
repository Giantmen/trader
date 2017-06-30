package config

type Server struct {
	Account   string //账户
	Accesskey string
	Secretkey string
	Timeout   int
}

type Config struct {
	Listen string

	Debug    bool
	LogPath  string
	LogLevel string

	Chbtc Server
	Yunbi Server
	Huobi Server
}
