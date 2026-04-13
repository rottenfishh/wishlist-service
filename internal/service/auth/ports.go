package auth

import (
	"cdek/internal/model"
	"context"
)

type Repository interface {
	SaveUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type Service interface {
	Register(ctx context.Context, email, password string) (*model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}
