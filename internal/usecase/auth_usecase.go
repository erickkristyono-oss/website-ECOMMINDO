package usecase

import (
	"ecommindo/internal/auth"
	"ecommindo/internal/domain"
	"ecommindo/internal/model/entity"
	"context"

	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	repo      domain.UserRepository
	jwtSecret []byte
}

func NewAuthUsecase(repo domain.UserRepository, jwtSecret []byte) domain.AuthUsecase {
	return &authUsecase{repo: repo, jwtSecret: jwtSecret}
}

func (u *authUsecase) Register(ctx context.Context, fullName, phone, email, password string) (*entity.User, string, error) {
	existing, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if existing != nil {
		return nil, "", ErrEmailTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &entity.User{
		FullName:     fullName,
		Phone:        phone,
		Email:        email,
		PasswordHash: string(hash),
	}

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := auth.GenerateToken(u.jwtSecret, user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (*entity.User, string, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(u.jwtSecret, user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (u *authUsecase) Me(ctx context.Context, userID int) (*entity.User, error) {
	return u.repo.FindByID(ctx, userID)
}
