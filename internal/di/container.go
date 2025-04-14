package di

import (
	au "github.com/seth16888/coauth/api/v1"
	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/config"
	"github.com/seth16888/wxbusiness/internal/data"
	"github.com/seth16888/wxbusiness/internal/handler"
	"github.com/seth16888/wxbusiness/internal/server"
	"github.com/seth16888/wxbusiness/pkg/jwt"
	"github.com/seth16888/wxbusiness/pkg/validator"
	ak "github.com/seth16888/wxtoken/api/v1"
	"go.uber.org/zap"
)

var di *Container

func init() {
	di = new(Container)
}

func Get() *Container {
	if di == nil {
		panic("di not initialized")
	}
	return di
}

type Container struct {
	Conf               *config.Conf // 配置文件
	DB                 *data.Data   // 数据库连接
	Log                *zap.Logger
	JWT                *jwt.JWTService
	Server             *server.Server
	HealthHandler      *handler.HealthHandler
	TokenClient        ak.TokenClient
	Validator          *validator.Validator
	PortalUsecase      *biz.PortalUsecase
	PlatformAppUsecase *biz.PlatformAppUsecase
	MenuUsecase        *biz.MPMenuUsecase
	CoAuthClient       au.CoauthClient
}
