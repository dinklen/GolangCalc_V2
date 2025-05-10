package models

import (
	"time"

	"github.com/google/uuid"
)

type AccountData struct {
	ID           uuid.UUID     `json:"id"` `db:"id"`
	Login        string        `json:"login"` `db:"login"`
	Password     string        `json:"password"` `db:"password"`
	CreatingTime time.Duration `json:"created_at"` `db:"created_at"`
}

