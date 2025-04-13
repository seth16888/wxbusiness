package bootstrap

import (
	"time"

	"github.com/seth16888/wxbusiness/internal/di"
	"github.com/seth16888/wxbusiness/internal/router"
	"github.com/seth16888/wxbusiness/internal/server"
	"github.com/seth16888/wxbusiness/pkg/jwt"
)

func InitServer(deps *di.Container) *server.Server {
  deps.JWT = jwt.NewJWTService(
    deps.Conf.Jwt.SignKey,
    deps.Conf.Jwt.Issuer,
    time.Duration(deps.Conf.Jwt.ExpireTime) * time.Second,
    time.Duration(deps.Conf.Jwt.MaxRefreshTime) * time.Second,
    time.Local,
    deps.Log,
  )

	r := router.InitRouter(deps)

  return server.NewServer(deps.Conf.Server, deps.Log, r)
}
