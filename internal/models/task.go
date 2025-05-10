package models

import "time"

type TaskStatus string

const (
    StatusPending   TaskStatus = "pending"
    StatusCompleted TaskStatus = "completed"
    StatusFailed    TaskStatus = "failed"
)

type Task struct {
    ID         string     `json:"id" db:"id"`
    Expression string     `json:"expression" db:"expression"`
    Result     *float64   `json:"result,omitempty" db:"result"`
    Status     TaskStatus `json:"status" db:"status"`
    Error      *string    `json:"error,omitempty" db:"error"`
    CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}
