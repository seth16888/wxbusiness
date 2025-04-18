package biz

import (
	"context"

	"github.com/seth16888/wxbusiness/internal/data/entities"
)

type UserAppRepo interface {
	GetByUserId(ctx context.Context, userId string) ([]*entities.PlatformApp, error)
}

type UserUsecase struct {
	repo UserAppRepo
}

func NewUserUsecase(repo UserAppRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) ListMPApps(ctx context.Context, userId string) ([]*entities.PlatformApp, error) {
	return u.repo.GetByUserId(ctx, userId)
}

