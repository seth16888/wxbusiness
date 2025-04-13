package data

import (
	"context"
	"fmt"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type PlatformAppData struct {
	col    *mongo.Collection
	data   *Data
	log    *zap.Logger
}

// Create adds a new Platform App
func (p *PlatformAppData) Create(ctx context.Context, app *entities.PlatformApp) (string, error) {
	// MongoDB会自动生成_id
	result, err := p.col.InsertOne(ctx, app)
	if err != nil {
		return "", err
	}

	// 将ObjectID转换为字符串返回
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	return "", fmt.Errorf("failed to get inserted id")
}

func NewPlatformAppRepo(data *Data, logger *zap.Logger) biz.PlatformAppRepo {
	collection := data.db.Collection("platform_apps")
	return &PlatformAppData{
		col:    collection,
		data:   data,
		log:    logger,
	}
}

// UpdateStatus Update Platform App status
func (p *PlatformAppData) UpdateStatus(ctx context.Context, id string, status int) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"status": status}}
	_, err = p.col.UpdateOne(ctx, filter, update)
	return err
}

// Get Platform App
func (p *PlatformAppData) Get(ctx context.Context, id string) (*entities.PlatformApp, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectID}
	app := &entities.PlatformApp{}
	if err := p.col.FindOne(ctx, filter).Decode(app); err != nil {
		return nil, err
	}
	return app, nil
}
