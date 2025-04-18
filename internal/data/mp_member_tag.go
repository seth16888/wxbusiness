package data

import (
	"context"
	"time"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MemberTagData struct {
	col  *mongo.Collection
	data *Data
	log  *zap.Logger
}

// Create implements biz.MemberTagRepo.
func (m *MemberTagData) Create(ctx context.Context, tag *entities.MemberTag) error {
	// mongoDB 插入
  now := time.Now().Unix()
  tag.CreatedAt = now
  tag.UpdatedAt = now
  _, err := m.col.InsertOne(ctx, tag)
  if err!= nil {
    return err
  }
  return nil
}

// Delete implements biz.MemberTagRepo.
func (m *MemberTagData) Delete(ctx context.Context, appId string, tagId int64) error {
  // mongoDB 删除
  filter := bson.M{"app_id": appId, "tag_id": tagId}
  _, err := m.col.DeleteOne(ctx, filter)
  if err!= nil {
    return err
  }
  return nil
}

// GetByTagId implements biz.MemberTagRepo.
func (m *MemberTagData) GetByTagId(ctx context.Context, appId string, tagId int64) (*entities.MemberTag, error) {
  // mongoDB 查询
  filter := bson.M{"app_id": appId, "tag_id": tagId}
  var tag entities.MemberTag
  err := m.col.FindOne(ctx, filter).Decode(&tag)
  if err!= nil {
    return nil, err
  }

  return &tag, nil
}

// Query implements biz.MemberTagRepo.
func (m *MemberTagData) Query(ctx context.Context, appId string) ([]*entities.MemberTag, error) {
  // mongoDB 查询
  filter := bson.M{"app_id": appId}
  cursor, err := m.col.Find(ctx, filter)
  if err!= nil {
    return nil, err
  }
  defer cursor.Close(ctx)
  var tags []*entities.MemberTag
  if err := cursor.All(ctx, &tags); err!= nil {
    return nil, err
  }
  return tags, nil
}

// Update implements biz.MemberTagRepo.
func (m *MemberTagData) Update(ctx context.Context, id primitive.ObjectID,  tag *entities.MemberTag) error {
  // mongoDB 更新
  now := time.Now().Unix()
  tag.UpdatedAt = now
  filter := bson.M{"_id": id}
  // 更新名称
  update := bson.M{"$set": bson.M{"name": tag.Name, "updated_at": now, "count": tag.Count}}
  _, err := m.col.UpdateOne(ctx, filter, update)
  if err!= nil {
    return err
  }
  return nil
}

// NewMemberTagData creates a new MemberTagData.
func NewMemberTagData(data *Data, log *zap.Logger) biz.MemberTagRepo {
	collection := data.db.Collection("member_tags")
	return &MemberTagData{data: data, log: log, col: collection}
}
