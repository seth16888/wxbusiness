package bootstrap

import (
	"github.com/seth16888/wxbusiness/internal/data"
	"github.com/seth16888/wxbusiness/internal/database"
	"github.com/seth16888/wxbusiness/internal/di"
	"go.uber.org/zap"
)

func InitDatabase(conf *database.DatabaseConfig, log *zap.Logger) error {
  log.Info("Initializing database connection...")
  // Initialize the database connection
	db := data.NewData(conf, log)

	di.Get().DB = db
	return nil
}
