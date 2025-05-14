package models

import (
	"time"

	"github.com/google/uuid"
)

type Expression struct {
	ID           uuid.UUID `json:"id"`
	Expr         string    `json:"expression"`
	Result       float64   `json:"result"`
	Status       string    `json:"status"`
	CreatingTime time.Time `json:"creating_time"`
}
