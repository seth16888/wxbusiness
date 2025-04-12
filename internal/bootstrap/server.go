package bootstrap

import (
	"github.com/seth16888/wxbusiness/internal/config"
	"github.com/seth16888/wxbusiness/internal/router"
	"github.com/seth16888/wxbusiness/internal/server"
	"go.uber.org/zap"
)

func InitServer(conf *config.Conf, log *zap.Logger) *server.Server {
	r := router.InitRouter(log)

	return server.NewServer(conf.Server, log, r)
}
