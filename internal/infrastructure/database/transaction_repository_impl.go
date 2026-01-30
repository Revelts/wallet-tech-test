package database

import (
	"context"
	"database/sql"
	"wallet-service/internal/domain/entity"
	"wallet-service/internal/domain/repository"
)

type TransactionRepositoryImpl struct {
	db *sql.DB
}

func NewTransactionRepositoryImpl(db *sql.DB) repository.TransactionRepository {
	return &TransactionRepositoryImpl{
		db: db,
	}
}

func (r *TransactionRepositoryImpl) getExecutor(ctx context.Context) executor {
	tx, ok := getTxFromContext(ctx)
	if ok {
		return tx
	}
	return r.db
}

func (r *TransactionRepositoryImpl) Create(ctx context.Context, transaction *entity.Transaction) (err error) {
	query := `INSERT INTO transactions (user_id, type, amount, created_at) VALUES (?, ?, ?, ?)`

	exec := r.getExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, transaction.UserID, transaction.Type, transaction.Amount, transaction.CreatedAt)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	transaction.ID = id
	return
}
