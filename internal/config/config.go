package config

import (
	"github.com/seth16888/wxbusiness/internal/database"
	"github.com/seth16888/wxbusiness/internal/server"
	"github.com/seth16888/wxbusiness/pkg/logger"

	"github.com/spf13/viper"
)

type Conf struct {
	Server      *server.ServerConf
	Log         *logger.LogConfig
	Redis       *RedisConfig
	DB          *database.DatabaseConfig
	Jwt         *Jwt
	TokenServer *TokenServer `mapstructure:"token_server"`
	ProxyServer *ProxyServer `mapstructure:"proxy_server"`
}

// TokenServer token server配置
type TokenServer struct {
	Addr string
}

// ProxyServer proxy server配置
type ProxyServer struct {
	Addr string
}

// redis配置
type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

type Jwt struct {
	Issuer         string
	ExpireTime     int
	MaxRefreshTime int
	SignKey        string
}

var appConf *Conf

func ReadConfigFromFile(file string) *Conf {
	defer func() {
		if r := recover(); r != nil {
			panic("Error reading config file: " + r.(error).Error() + "\n")
		}
	}()

	if file == "" {
		file = "conf.yaml"
	}

	viper.SetConfigFile(file)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath("~")

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	appConf = &Conf{}
	if err := viper.Unmarshal(appConf); err != nil {
		panic(err)
	}

	// watch
	viper.WatchConfig()

	return appConf
}
