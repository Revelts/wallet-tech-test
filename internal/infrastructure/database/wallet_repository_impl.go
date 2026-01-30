package database

import (
	"context"
	"database/sql"
	"errors"
	"wallet-service/internal/domain/entity"
	"wallet-service/internal/domain/repository"
)

type WalletRepositoryImpl struct {
	db *sql.DB
}

func NewWalletRepositoryImpl(db *sql.DB) repository.WalletRepository {
	return &WalletRepositoryImpl{
		db: db,
	}
}

type executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func (r *WalletRepositoryImpl) getExecutor(ctx context.Context) executor {
	tx, ok := getTxFromContext(ctx)
	if ok {
		return tx
	}
	return r.db
}

func (r *WalletRepositoryImpl) Create(ctx context.Context, wallet *entity.Wallet) (err error) {
	query := `INSERT INTO wallets (user_id, balance, created_at) VALUES (?, ?, ?)`

	exec := r.getExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, wallet.UserID, wallet.Balance, wallet.CreatedAt)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	wallet.ID = id
	return
}

func (r *WalletRepositoryImpl) FindByUserID(ctx context.Context, userID int64) (wallet *entity.Wallet, err error) {
	query := `SELECT id, user_id, balance, created_at FROM wallets WHERE user_id = ?`

	exec := r.getExecutor(ctx)
	wallet = &entity.Wallet{}
	err = exec.QueryRowContext(ctx, query, userID).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.CreatedAt,
	)

	if err == sql.ErrNoRows {
		err = errors.New("wallet not found")
		return
	}

	if err != nil {
		return
	}

	return
}

func (r *WalletRepositoryImpl) UpdateBalance(ctx context.Context, userID int64, newBalance int64) (err error) {
	query := `UPDATE wallets SET balance = ? WHERE user_id = ?`

	exec := r.getExecutor(ctx)
	result, err := exec.ExecContext(ctx, query, newBalance, userID)
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = errors.New("wallet not found")
		return
	}

	return
}
