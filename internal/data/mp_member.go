package data

import (
	"context"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/data/entities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MPMemberData struct {
	col  *mongo.Collection
	data *Data
	log  *zap.Logger
}

// FindByAppId implements biz.MPMemberRepo.
func (m *MPMemberData) FindByAppId(c context.Context, appId string) ([]*entities.MPMember, error) {
	panic("unimplemented")
}

// FindById implements biz.MPMemberRepo.
func (m *MPMemberData) FindById(c context.Context, id string) (*entities.MPMember, error) {
	panic("unimplemented")
}

// UpdateRemark implements biz.MPMemberRepo.
func (m *MPMemberData) UpdateRemark(c context.Context, id string, remark string) error {
	panic("unimplemented")
}

// NewMPMemberData returns a new MPMemberData.
func NewMPMemberData(data *Data, log *zap.Logger) biz.MPMemberRepo {
	collection := data.db.Collection("mp_members")
	return &MPMemberData{col: collection, data: data, log: log}
}
