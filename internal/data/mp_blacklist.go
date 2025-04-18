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

type MPBlackListData struct {
	col  *mongo.Collection
	data *Data
	log  *zap.Logger
}

// Block 批量封禁
func (m *MPBlackListData) Block(c context.Context, appId string, openids []string) error {
	// 查询
  filter:= bson.M{
    "app_id": appId,
    "openid": bson.M{"$in": openids},
  }
  update := bson.M{
    "$set": bson.M{
      "blocked": true,
    },
  }

  _, err := m.col.UpdateMany(c, filter, update)
  return err
}

// Query 查询黑名单
func (m *MPBlackListData) Query(c context.Context, appId string) ([]*entities.MPMember, error) {
  filter := bson.M{
    "app_id": appId,
    "blocked": true,
  }
  opts := options.Find().SetSort(bson.M{"created_at": -1})
  cursor, err := m.col.Find(c, filter, opts)
  if err != nil {
    return nil, err
  }
  defer cursor.Close(c)

  var blacklists []*entities.MPMember
  if err = cursor.All(c, &blacklists); err != nil {
    return nil, err
  }

  return blacklists, nil
}

// Unblock 批量解封
func (m *MPBlackListData) Unblock(c context.Context, appId string, openids []string) error {
	// 查询
  filter:= bson.M{
    "app_id": appId,
    "openid": bson.M{"$in": openids},
  }
  update := bson.M{
    "$set": bson.M{
      "blocked": false,
    },
  }

  _, err := m.col.UpdateMany(c, filter, update)
  return err
}

// NewMPBlackListData
func NewMPBlackListData(data *Data, log *zap.Logger) biz.MPBlackListRepo {
	collection := data.db.Collection("mp_members")
	return &MPBlackListData{log: log, data: data, col: collection}
}
