package user

import (
	"context"
	"gateway/internal/model"
)

type Service interface {
	Create(ctx context.Context, username, password, email string) error
	GetById(ctx context.Context, id uint64) (model.User, error)
	Update(ctx context.Context, id uint64, username, password, email string) error
	Login(ctx context.Context, username, password string) (string, error)
}

type Controller struct {
	service Service
}
