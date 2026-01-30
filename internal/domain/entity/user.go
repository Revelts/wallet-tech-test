package entity

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	PinHash      string
	CreatedAt    time.Time
}
