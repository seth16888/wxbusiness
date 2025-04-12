package bootstrap

import "github.com/seth16888/wxbusiness/internal/config"


func InitConfig(configFile string) (conf *config.Conf, err error) {
	return config.ReadConfigFromFile(configFile), nil
}
