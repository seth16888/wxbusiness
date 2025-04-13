package data

import (
	"context"
	"time"

	"github.com/seth16888/wxbusiness/internal/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Data struct {
	client *mongo.Client
	log    *zap.Logger
	dbConf *database.DatabaseConfig
	db     *mongo.Database
}

// NewData new data
func NewData(conf *database.DatabaseConfig, log *zap.Logger) *Data {
	client := newDB(conf, log)
	db := client.Database(conf.DatabaseName)
	return &Data{
		client: client,
		log:    log,
		dbConf: conf,
		db:     db,
	}
}

// NewDB
func newDB(conf *database.DatabaseConfig, log *zap.Logger) *mongo.Client {
	clientOptions := options.Client().ApplyURI(conf.Source).
		SetMaxPoolSize(conf.MaxPoolSize).
		SetMinPoolSize(conf.MinPoolSize).
		SetMaxConnIdleTime(time.Duration(conf.MaxIdleTime) * time.Second)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("failed connecting to MongoDB", zap.Error(err))
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("failed pinging MongoDB", zap.Error(err))
	}

	log.Info("Connected to MongoDB")
	return client
}
