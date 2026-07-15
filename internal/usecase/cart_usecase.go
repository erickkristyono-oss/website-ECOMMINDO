package usecase

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"context"
)

type cartUsecase struct {
	cartRepo    domain.CartRepository
	serviceRepo domain.ServiceRepository
}

func NewCartUsecase(cartRepo domain.CartRepository, serviceRepo domain.ServiceRepository) domain.CartUsecase {
	return &cartUsecase{cartRepo: cartRepo, serviceRepo: serviceRepo}
}

func (u *cartUsecase) Add(ctx context.Context, userID int, serviceID string) (*domain.CartSummary, error) {
	if _, err := u.serviceRepo.FindByID(ctx, serviceID); err != nil {
		return nil, ErrServiceNotFound
	}

	if err := u.cartRepo.Add(ctx, userID, serviceID); err != nil {
		return nil, err
	}

	return u.Get(ctx, userID)
}

func (u *cartUsecase) Remove(ctx context.Context, userID int, serviceID string) (*domain.CartSummary, error) {
	if err := u.cartRepo.Remove(ctx, userID, serviceID); err != nil {
		return nil, err
	}

	return u.Get(ctx, userID)
}

func (u *cartUsecase) Get(ctx context.Context, userID int) (*domain.CartSummary, error) {
	items, err := u.cartRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	summary := &domain.CartSummary{Items: []entity.Service{}}
	for _, item := range items {
		service, err := u.serviceRepo.FindByID(ctx, item.ServiceID)
		if err != nil {
			continue
		}
		summary.Items = append(summary.Items, *service)
		summary.Total += service.Price
	}

	return summary, nil
}
