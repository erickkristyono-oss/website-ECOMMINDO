package domain

import (
	"ecommindo/internal/model/entity"
	"context"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	FindByUser(ctx context.Context, userID int) ([]entity.Order, error)
}

type OrderUsecase interface {
	Checkout(ctx context.Context, userID int) (*entity.Order, error)
	ListByUser(ctx context.Context, userID int) ([]entity.Order, error)
}
