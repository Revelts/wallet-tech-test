package database

import (
	"context"
	"database/sql"
	"wallet-service/internal/domain/entity"
	"wallet-service/internal/domain/repository"
)

type APILogRepositoryImpl struct {
	db *sql.DB
}

func NewAPILogRepositoryImpl(db *sql.DB) repository.APILogRepository {
	return &APILogRepositoryImpl{
		db: db,
	}
}

func (r *APILogRepositoryImpl) Create(ctx context.Context, log *entity.APILog) (err error) {
	query := `INSERT INTO api_logs 
		(request_id, timestamp, http_method, request_path, http_status, latency_ms, 
		ip_address, user_agent, user_id, request_body, response_body, created_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var userIDValue interface{}
	if log.UserID > 0 {
		userIDValue = log.UserID
	} else {
		userIDValue = nil
	}

	result, err := r.db.ExecContext(ctx, query,
		log.RequestID,
		log.Timestamp,
		log.HTTPMethod,
		log.RequestPath,
		log.HTTPStatus,
		log.LatencyMs,
		log.IPAddress,
		log.UserAgent,
		userIDValue,
		log.RequestBody,
		log.ResponseBody,
		log.CreatedAt,
	)
	if err != nil {
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	log.ID = id
	return
}
