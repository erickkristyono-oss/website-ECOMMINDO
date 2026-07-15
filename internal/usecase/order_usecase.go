package usecase

import (
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"context"
	"fmt"
	"math/rand"
	"time"
)

type orderUsecase struct {
	orderRepo domain.OrderRepository
	cartUC    domain.CartUsecase
	cartRepo  domain.CartRepository
}

func NewOrderUsecase(orderRepo domain.OrderRepository, cartUC domain.CartUsecase, cartRepo domain.CartRepository) domain.OrderUsecase {
	return &orderUsecase{orderRepo: orderRepo, cartUC: cartUC, cartRepo: cartRepo}
}

func (u *orderUsecase) Checkout(ctx context.Context, userID int) (*entity.Order, error) {
	cart, err := u.cartUC.Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, ErrCartEmpty
	}

	order := &entity.Order{
		OrderCode: generateOrderCode(),
		UserID:    userID,
		Total:     cart.Total,
		Status:    "Menunggu Konfirmasi Pembayaran",
	}

	for _, s := range cart.Items {
		order.Items = append(order.Items, entity.OrderItem{
			ServiceID:   s.ID,
			ServiceName: s.Name,
			Price:       s.Price,
		})
	}

	if err := u.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	if err := u.cartRepo.Clear(ctx, userID); err != nil {
		return nil, err
	}

	return order, nil
}

func (u *orderUsecase) ListByUser(ctx context.Context, userID int) ([]entity.Order, error) {
	return u.orderRepo.FindByUser(ctx, userID)
}

func generateOrderCode() string {
	return fmt.Sprintf("EJP-%d%03d", time.Now().Unix()%1000000, rand.Intn(1000))
}
