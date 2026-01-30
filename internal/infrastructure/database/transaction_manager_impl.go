package database

import (
	"context"
	"database/sql"
	"wallet-service/internal/domain/repository"
)

type TransactionManagerImpl struct {
	db *sql.DB
}

func NewTransactionManagerImpl(db *sql.DB) repository.TransactionManager {
	return &TransactionManagerImpl{
		db: db,
	}
}

type txKey struct{}

func (tm *TransactionManagerImpl) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := tm.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	txCtx := context.WithValue(ctx, txKey{}, tx)
	err = fn(txCtx)

	return
}

func getTxFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}
