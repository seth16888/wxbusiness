package data

import (
	"context"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MPMaterialData struct {
	col  *mongo.Collection
	data *Data
	log  *zap.Logger
}

// SaveMany implements biz.MaterialRepo.
func (m *MPMaterialData) SaveMany(c context.Context, materials []*entities.MPMaterial) error {
	// mongoDB insert many, if not exist, insert it.
  for _, material := range materials {
    filter := bson.M{
      "app_id":       material.AppId,
      "media_id":     material.MediaId,
    }
    update := bson.M{
      "$set": material,
    }
    _, err := m.col.UpdateOne(c, filter, update, options.Update().SetUpsert(true))
    if err!= nil {
      return err
    }
  }
  return nil
}

// Find implements biz.MaterialRepo.
func (m *MPMaterialData) Find(c context.Context, appId string,
	IsPermanent bool, mediaType string, pageNo int64, pageSize int64,
) ([]*entities.MPMaterial, error) {
	// mongoDB find
	filter := bson.M{
		"app_id":       appId,
		"media_type":   mediaType,
		"is_permanent": IsPermanent,
	}
	skip := GetSkipNum(pageNo, pageSize)
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(skip).
		SetLimit(pageSize)
	cursor, err := m.col.Find(c, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var materials []*entities.MPMaterial
	if err = cursor.All(c, &materials); err != nil {
		return nil, err
	}
	return materials, nil
}

// Insert implements biz.MaterialRepo.
func (m *MPMaterialData) Insert(c context.Context,
	material *entities.MPMaterial,
) error {
	// mongoDB insert
	_, err := m.col.InsertOne(c, material)
	return err
}

// NewMPMaterialData
func NewMPMaterialData(data *Data, log *zap.Logger) biz.MaterialRepo {
	collection := data.db.Collection("mp_materials")
	return &MPMaterialData{col: collection, data: data, log: log}
}
