package repository

import (
	"context"
	"wallet-service/internal/domain/entity"
)

type APILogRepository interface {
	Create(ctx context.Context, log *entity.APILog) (err error)
}
