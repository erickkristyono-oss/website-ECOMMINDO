package domain

import (
	"ecommindo/internal/model/entity"
	"context"
)

type ServiceRepository interface {
	FindAll(ctx context.Context) ([]entity.Service, error)
	FindByID(ctx context.Context, id string) (*entity.Service, error)
}

type ServiceUsecase interface {
	GetAll(ctx context.Context) ([]entity.Service, error)
}
