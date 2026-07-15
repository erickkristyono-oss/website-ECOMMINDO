package domain

import (
	"ecommindo/internal/model/entity"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
}

type AuthUsecase interface {
	Register(ctx context.Context, fullName, phone, email, password string) (*entity.User, string, error)
	Login(ctx context.Context, email, password string) (*entity.User, string, error)
	Me(ctx context.Context, userID int) (*entity.User, error)
}
