package domain

import (
	"ecommindo/internal/model/entity"
	"context"
)

type CartRepository interface {
	Add(ctx context.Context, userID int, serviceID string) error
	Remove(ctx context.Context, userID int, serviceID string) error
	FindByUser(ctx context.Context, userID int) ([]entity.CartItem, error)
	Clear(ctx context.Context, userID int) error
}

type CartSummary struct {
	Items []entity.Service `json:"items"`
	Total float64          `json:"total"`
}

type CartUsecase interface {
	Add(ctx context.Context, userID int, serviceID string) (*CartSummary, error)
	Remove(ctx context.Context, userID int, serviceID string) (*CartSummary, error)
	Get(ctx context.Context, userID int) (*CartSummary, error)
}
