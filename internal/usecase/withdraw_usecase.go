package usecase

import (
	"context"
	"time"
	"wallet-service/internal/app_error"
	"wallet-service/internal/domain/entity"
	"wallet-service/internal/domain/repository"

	"golang.org/x/crypto/bcrypt"
)

type WithdrawUsecase struct {
	userRepo        repository.UserRepository
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	txManager       repository.TransactionManager
}

func NewWithdrawUsecase(
	userRepo repository.UserRepository,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	txManager repository.TransactionManager,
) *WithdrawUsecase {
	return &WithdrawUsecase{
		userRepo:        userRepo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		txManager:       txManager,
	}
}

type WithdrawInput struct {
	UserID int64
	Amount int64
	Pin    string
}

type WithdrawOutput struct {
	TransactionID int64
	NewBalance    int64
}

func (u *WithdrawUsecase) Execute(ctx context.Context, input WithdrawInput) (output WithdrawOutput, err error) {
	if input.Amount <= 0 {
		err = app_error.InvalidAmount
		return
	}

	user, err := u.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PinHash), []byte(input.Pin))
	if err != nil {
		err = app_error.InvalidPIN
		return
	}

	wallet, err := u.walletRepo.FindByUserID(ctx, input.UserID)
	if err != nil {
		return
	}

	if wallet.Balance < input.Amount {
		err = app_error.InsufficientBalance
		return
	}

	newBalance := wallet.Balance - input.Amount

	var transactionID int64

	err = u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		updateErr := u.walletRepo.UpdateBalance(txCtx, input.UserID, newBalance)
		if updateErr != nil {
			return updateErr
		}

		transaction := &entity.Transaction{
			UserID:    input.UserID,
			Type:      entity.TransactionTypeDebit,
			Amount:    input.Amount,
			CreatedAt: time.Now(),
		}

		createErr := u.transactionRepo.Create(txCtx, transaction)
		if createErr != nil {
			return createErr
		}

		transactionID = transaction.ID
		return nil
	})

	if err != nil {
		return
	}

	output = WithdrawOutput{
		TransactionID: transactionID,
		NewBalance:    newBalance,
	}

	return
}
