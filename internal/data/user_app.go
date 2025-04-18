package data

import (
	"context"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserAppData struct {
  col  *mongo.Collection
	data *Data
	log  *zap.Logger
}

// GetByUserId
func (u *UserAppData) GetByUserId(ctx context.Context, userId string) ([]*entities.PlatformApp, error) {
	filter := bson.M{"user_id": userId}
	cursor, err := u.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var apps []*entities.PlatformApp
	if err := cursor.All(ctx, &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

// NewUserAppData creates a new UserAppData.
func NewUserAppData(data *Data, log *zap.Logger) biz.UserAppRepo {
	collection := data.db.Collection("platform_apps")
	return &UserAppData{data: data, log: log, col: collection}
}
