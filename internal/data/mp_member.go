package data

import (
	"context"
	"time"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/data/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MPMemberData struct {
	col  *mongo.Collection
	data *Data
	log  *zap.Logger
}

// Save implements biz.MPMemberRepo.
func (m *MPMemberData) Save(c context.Context, members []*entities.MPMember) error {
	// 查询: app_id, openid
  now:= time.Now().Unix()
	for _, member := range members {
		filter := bson.M{"app_id": member.AppId, "openid": member.OpenId}
		update := bson.M{"$set": bson.M{
			"mp_id":           member.MpId,
			"subscribe":       member.Subscribe,
			"nick_name":       member.NickName,
			"sex":             member.Sex,
			"language":        member.Language,
			"city":            member.City,
			"province":        member.Province,
			"country":         member.Country,
			"subscribe_time":  member.SubscribeTime,
			"union_id":        member.UnionId,
			"remark":          member.Remark,
			"group_id":        member.GroupId,
			"tags":            member.Tags,
			"subscribe_scene": member.SubscribeScene,
			"qr_scene":        member.QrScene,
			"qr_scene_str":    member.QrSceneStr,
			"message_count":   member.MessageCount,
			"comment_count":   member.CommentCount,
			"star_comment":    member.StarComment,
			"praise_count":    member.PraiseCount,
			"praise_amounts":  member.PraiseAmounts,
			"updated_at":      time.Now().Unix(), // 只更新更新时间
			"blocked":         member.Blocked,
		}}
		result, err := m.col.UpdateOne(c, filter, update)
		if err != nil {
			return err
		}
		if result.MatchedCount == 0 {
			// 不存在则创建
			member.CreatedAt = now
			member.UpdatedAt = now
			_, err := m.col.InsertOne(c, member)
			if err!= nil {
				return err
			}
		}
	}
  return nil
}

// FindByAppId implements biz.MPMemberRepo.
func (m *MPMemberData) FindByAppId(c context.Context, appId string) ([]*entities.MPMember, error) {
	filter := bson.M{"app_id": appId}
	cursor, err := m.col.Find(c, filter)
	if err != nil {
		m.log.Error("find member by app id error", zap.Error(err))
		return nil, err
	}
	var members []*entities.MPMember
	if err := cursor.All(c, &members); err != nil {
		m.log.Error("decode member error", zap.Error(err))
		return nil, err
	}

	return members, nil
}

// FindById implements biz.MPMemberRepo.
func (m *MPMemberData) FindById(c context.Context, id string) (*entities.MPMember, error) {
	filter := bson.M{"_id": id}
	var member entities.MPMember
	err := m.col.FindOne(c, filter).Decode(&member)
	if err != nil {
		m.log.Error("find member by id error", zap.Error(err))
		return nil, err
	}

	return &member, nil
}

// UpdateRemark implements biz.MPMemberRepo.
func (m *MPMemberData) UpdateRemark(c context.Context, id string, remark string) error {
	// mongoDB update
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"remark": remark}}
	_, err := m.col.UpdateOne(c, filter, update)
	if err != nil {
		m.log.Error("update remark error", zap.Error(err))
		return err
	}

	return nil
}

// NewMPMemberData returns a new MPMemberData.
func NewMPMemberData(data *Data, log *zap.Logger) biz.MPMemberRepo {
	collection := data.db.Collection("mp_members")
	return &MPMemberData{col: collection, data: data, log: log}
}
