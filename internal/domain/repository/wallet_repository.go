package repository

import (
	"context"
	"wallet-service/internal/domain/entity"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *entity.Wallet) (err error)
	FindByUserID(ctx context.Context, userID int64) (wallet *entity.Wallet, err error)
	UpdateBalance(ctx context.Context, userID int64, newBalance int64) (err error)
}
