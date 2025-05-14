package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountData struct {
	ID           uuid.UUID `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"password"`
	CreatingTime time.Time `json:"created_at"`
}
