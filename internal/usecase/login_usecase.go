package usecase

import (
	"context"
	"wallet-service/internal/app_error"
	"wallet-service/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase struct {
	userRepo repository.UserRepository
}

func NewLoginUsecase(userRepo repository.UserRepository) *LoginUsecase {
	return &LoginUsecase{
		userRepo: userRepo,
	}
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	UserID int64
	Email  string
}

func (u *LoginUsecase) Execute(ctx context.Context, input LoginInput) (output LoginOutput, err error) {
	user, err := u.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		err = app_error.InvalidCredentials
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		err = app_error.InvalidCredentials
		return
	}

	output = LoginOutput{
		UserID: user.ID,
		Email:  user.Email,
	}

	return
}
