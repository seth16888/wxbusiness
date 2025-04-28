package data

import (
	"context"

	"github.com/seth16888/wxbusiness/internal/biz"
	"github.com/seth16888/wxbusiness/internal/data/entities"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type MPMenuData struct {
	log  *zap.Logger
	col  *mongo.Collection
	data *Data
}

// DeleteMenu implements biz.MenuRepo.
func (m *MPMenuData) DeleteMenu(ctx context.Context, pId string) error {
	filter := map[string]interface{}{
		"app_id": pId,
	}
	_, err := m.col.DeleteOne(ctx, filter)
	return err
}

// GetMenuInfo implements biz.MenuRepo.
func (m *MPMenuData) GetMenuInfo(ctx context.Context, appId string) (*entities.MPMenu, error) {
	filter := map[string]interface{}{
		"app_id": appId,
	}
	var menuRes entities.MPMenu
	err := m.col.FindOne(ctx, filter).Decode(&menuRes)
	return &menuRes, err
}

// SaveMenu implements biz.MenuRepo.
func (m *MPMenuData) SaveMenu(ctx context.Context, apiMenu *entities.MPMenu) error {
	_, err := m.col.InsertOne(ctx, apiMenu)
	return err
}

// NewMPMenuData creates a new MPMenuData
func NewMPMenuData(log *zap.Logger, data *Data) biz.MenuRepo {
	return &MPMenuData{
		log:  log,
		col:  data.db.Collection("mp_menus"),
		data: data,
	}
}
