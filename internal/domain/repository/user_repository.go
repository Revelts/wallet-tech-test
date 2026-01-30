package repository

import (
	"context"
	"wallet-service/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (err error)
	FindByEmail(ctx context.Context, email string) (user *entity.User, err error)
	FindByID(ctx context.Context, id int64) (user *entity.User, err error)
}
