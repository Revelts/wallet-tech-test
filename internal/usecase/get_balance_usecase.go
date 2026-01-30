package usecase

import (
	"context"
	"wallet-service/internal/domain/repository"
)

type GetBalanceUsecase struct {
	walletRepo repository.WalletRepository
}

func NewGetBalanceUsecase(walletRepo repository.WalletRepository) *GetBalanceUsecase {
	return &GetBalanceUsecase{
		walletRepo: walletRepo,
	}
}

type GetBalanceInput struct {
	UserID int64
}

type GetBalanceOutput struct {
	Balance int64
}

func (u *GetBalanceUsecase) Execute(ctx context.Context, input GetBalanceInput) (output GetBalanceOutput, err error) {
	wallet, err := u.walletRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return
	}

	output = GetBalanceOutput{
		Balance: wallet.Balance,
	}

	return
}
