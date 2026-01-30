package usecase

import (
	"context"
	"time"
	"wallet-service/internal/app_error"
	"wallet-service/internal/domain/entity"
	"wallet-service/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUsecase struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
}

func NewRegisterUsecase(userRepo repository.UserRepository, walletRepo repository.WalletRepository) *RegisterUsecase {
	return &RegisterUsecase{
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

type RegisterInput struct {
	Email    string
	Password string
	Pin      string
}

type RegisterOutput struct {
	UserID int64
	Email  string
}

func (u *RegisterUsecase) Execute(ctx context.Context, input RegisterInput) (output RegisterOutput, err error) {
	existingUser, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		err = app_error.EmailAlreadyExists
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	pinHash, err := bcrypt.GenerateFromPassword([]byte(input.Pin), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	user := &entity.User{
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		PinHash:      string(pinHash),
		CreatedAt:    time.Now(),
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return
	}

	wallet := &entity.Wallet{
		UserID:    user.ID,
		Balance:   0,
		CreatedAt: time.Now(),
	}

	err = u.walletRepo.Create(ctx, wallet)
	if err != nil {
		return
	}

	output = RegisterOutput{
		UserID: user.ID,
		Email:  user.Email,
	}

	return
}
