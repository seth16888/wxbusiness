package database

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// NewDB 提供数据库连接
func NewDB(cfg *DatabaseConfig) (*gorm.DB, error) {
	if cfg.Source == "" {
		return nil, fmt.Errorf("database source is required")
	}

	if cfg.Driver != "mysql" {
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	logger := NewGormLogger(cfg.LogLevel)

	// 建立连接
	dial := mysql.New(mysql.Config{
		DriverName:                "mysql",
		DSN:                       cfg.Source,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})

	db, err := gorm.Open(dial, &gorm.Config{
		Logger: logger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
		NamingStrategy: schema.NamingStrategy{
			SingularTable: cfg.SingularTable,
			TablePrefix:   cfg.TablePrefix,
		},
		PrepareStmt:       cfg.PrepareStmt,
		AllowGlobalUpdate: cfg.AllowGlobalUpdate,
	})
	if err != nil {
		zap.L().Error("数据库连接失败", zap.Error(err), zap.String("DSN", cfg.Source))
		return nil, err
	}

	// 连接池设置
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("获取数据库连接池失败", zap.Error(err))
		return nil, err
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	zap.L().Info("数据库连接成功",
		zap.String("DSN", cfg.Source))
	return db, nil
}
