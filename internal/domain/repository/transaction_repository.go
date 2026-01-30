package repository

import (
	"context"
	"wallet-service/internal/domain/entity"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) (err error)
}
