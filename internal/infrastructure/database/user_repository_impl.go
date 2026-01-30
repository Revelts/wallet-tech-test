package database

import (
	"context"
	"database/sql"
	"errors"
	"wallet-service/internal/domain/entity"
	"wallet-service/internal/domain/repository"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepositoryImpl(db *sql.DB) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) (err error) {
	query := `INSERT INTO users (email, password_hash, pin_hash, created_at) VALUES (?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query, user.Email, user.PasswordHash, user.PinHash, user.CreatedAt)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	user.ID = id
	return
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (user *entity.User, err error) {
	query := `SELECT id, email, password_hash, pin_hash, created_at FROM users WHERE email = ?`

	user = &entity.User{}
	err = r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.PinHash,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		err = errors.New("user not found")
		return
	}

	if err != nil {
		return
	}

	return
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id int64) (user *entity.User, err error) {
	query := `SELECT id, email, password_hash, pin_hash, created_at FROM users WHERE id = ?`

	user = &entity.User{}
	err = r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.PinHash,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		err = errors.New("user not found")
		return
	}

	if err != nil {
		return
	}

	return
}
