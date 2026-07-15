package usecase

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"context"
)

type serviceUsecase struct {
	repo domain.ServiceRepository
}

func NewServiceUsecase(repo domain.ServiceRepository) domain.ServiceUsecase {
	return &serviceUsecase{repo: repo}
}

func (u *serviceUsecase) GetAll(ctx context.Context) ([]entity.Service, error) {
	return u.repo.FindAll(ctx)
}
