package di

import (
	"github.com/seth16888/wxbusiness/internal/config"
	"github.com/seth16888/wxbusiness/internal/database"
	"github.com/seth16888/wxbusiness/internal/handler"
	"github.com/seth16888/wxbusiness/internal/server"
	"github.com/seth16888/wxbusiness/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	Conf          *config.Conf // 配置文件
	DB            *gorm.DB     // 数据库连接
	Log           *zap.Logger
	Server        *server.Server
	HealthHandler *handler.HealthHandler
}

func NewContainer(configFile string) *Container {
	conf := config.ReadConfigFromFile(configFile)
	log := logger.InitLogger(conf.Log)

	db, err := database.NewDB(conf.DB)
	if err != nil {
		log.Fatal("failed to connect database", zap.Error(err))
	}

	healthHandler := handler.NewHealthHandler()

	di = &Container{
		Conf:          conf,
		DB:            db,
		Log:           log,
		HealthHandler: healthHandler,
	}
	return di
}
